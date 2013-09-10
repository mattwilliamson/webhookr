import logging
import os
from socketio import socketio_manage
from gevent import monkey
from flask import Flask, Response, request, render_template, url_for, redirect
from socketio.virtsocket import Socket
from socketio.server import SocketIOServer
from sockets import WebhookNamespace
from utils import formatted_json, scalar_or_list
import shortid

logging.basicConfig(level=logging.DEBUG)



app = Flask(__name__)
app.debug = os.getenv('ENVIRONMENT') == 'DEVELOPMENT'
app.config.from_object('default_settings')

# Allow a cfg to override default settings
if 'WEBHOOKR_SETTINGS' in os.environ:
    app.config.from_envvar('WEBHOOKR_SETTINGS')

# Patch up IO for gevent
monkey.patch_all()


@app.route('/')
def home():
    return render_template('webhooks/index.html')

@app.route('/help')
def help():
    webhook_id = shortid.generate()
    return render_template('webhooks/help.html', webhook_id=webhook_id, supported_sounds=app.config['VALID_SOUNDS'])


@app.route('/new')
def new():
    webhook_id = shortid.generate()
    return redirect(url_for('webhook', webhook_id=webhook_id))


@app.route('/<webhook_id>', methods=['GET', 'POST', 'HEAD', 'PUT', 'DELETE', 'OPTIONS'])
@app.route('/<webhook_id>/<path:remaining>', methods=['GET', 'POST', 'HEAD', 'PUT', 'DELETE', 'OPTIONS'])
def webhook(webhook_id, remaining=''):
    """Webhook callback requests this or when viewing a webhook"""

    has_data = len(request.args) or request.method != 'GET' or remaining
    post = formatted_json(request.form) if request.form else scalar_or_list(request.data)
    get = formatted_json(request.args)
    headers = formatted_json({k: v for k, v in request.headers.iteritems()})
    files = formatted_json({file_key: file.filename for file_key, file in request.files.iteritems()})

    # Play sound effect if url is like http://webhookr.com/aBr3/loud or /loud/something/else/
    sound_effect = app.config['VALID_SOUNDS'][0]
    for sound in app.config['VALID_SOUNDS']:
        sound_prefix = '{}/'.format(sound)
        if remaining == sound or remaining.startswith(sound_prefix):
            sound_effect = sound

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
        "hasData": has_data,
        "soundEffect": sound_effect,
    }

    if has_data:
        # Little hack to communicate to server outside of a websocket
        request.environ['socketio'] = Socket(server=app.server)
        wns = WebhookNamespace(request.environ, '/webhooks', request)

        wns.emit_to_channel(webhook_id, 'new_request', data)
        return "Message Logged"

    context = {'webhook_id': webhook_id, "full_path": url_for('webhook', webhook_id=webhook_id, _external=True)}

    return render_template('webhooks/webhook.html', **context)



@app.route('/socket.io/<path:remaining>')
def socketio(remaining):
    try:
        socketio_manage(request.environ, {'/webhooks': WebhookNamespace}, request)
    except:
        app.logger.error("Exception while handling socketio connection", exc_info=True)

    return Response()


if __name__ == '__main__':
    http_binding = (app.config['HOST'], app.config['PORT'])

    print 'app.debug = %s' % app.debug
    print 'Listening on %s:%s for HTTP...' % http_binding
    print 'Listening on %s:10843 for Flash policy server...' % http_binding[0]

    app.server = SocketIOServer(http_binding, app, policy_server=True, resource="socket.io")
    app.server.serve_forever()