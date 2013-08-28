import uuid
from socketio import socketio_manage
from gevent import monkey
from flask import Flask, Response, request, render_template, url_for, redirect
from socketio.virtsocket import Socket
from sockets import WebhookNamespace
from utils import nice_dict, scalar_or_list


app = Flask(__name__)
app.debug = True
app.config['MAX_CONTENT_LENGTH'] = 500 * 1024 # 500kb

monkey.patch_all()


@app.route('/')
def home():
    return render_template('webhooks/index.html')


@app.route('/new')
def new():
    webhook_id = str(uuid.uuid4())
    return redirect(url_for('webhook', webhook_id=webhook_id))


@app.route('/<webhook_id>', methods=['GET', 'POST', 'HEAD', 'PUT', 'DELETE', 'OPTIONS'])
@app.route('/<webhook_id>/<path:remaining>', methods=['GET', 'POST', 'HEAD', 'PUT', 'DELETE', 'OPTIONS'])
def webhook(webhook_id, remaining=None):
    if len(request.args) or request.method != 'GET':
        post = nice_dict(request.form) if request.form else scalar_or_list(request.data)
        get = nice_dict(request.args)
        headers = nice_dict({k: v for k, v in request.headers.iteritems()})
        files = nice_dict({file_key: file.filename for file_key, file in request.files.iteritems()})
        data = {
            "method": request.method,
            "ip": request.remote_addr,
            "host": request.headers.get('Host'),
            "length": request.headers.get('Content-Length'),
            "type": request.headers.get('Content-Type'),
            "post": post,
            "get": get,
            "files": files,
            "headers": headers,
            "path": remaining,
        }

        # Little hack to communicate to server outside of a websocket
        request.environ['socketio'] = Socket(server=app.server)
        wns = WebhookNamespace(request.environ, '/webhooks', request)
        wns.emit_to_room(webhook_id, 'new_request', data)

        return "Message Logged"

    context = {'webhook_id': webhook_id, "full_path": url_for('webhook', webhook_id=webhook_id, _external=True)}

    return render_template('webhooks/webhook.html', **context)



@app.route('/socket.io/<path:remaining>')
def socketio(remaining):
    try:
        socketio_manage(request.environ, {'/webhooks': WebhookNamespace}, request)
    except:
        app.logger.error("Exception while handling socketio connection",
                         exc_info=True)
    return Response()


if __name__ == '__main__':
    app.run()