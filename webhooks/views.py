import uuid

from utils import nice_dict, scalar_or_list

from django.shortcuts import render, redirect
from django.core.urlresolvers import reverse
from django.http import HttpResponse
from django.views.decorators.csrf import csrf_exempt
from django_socketio import broadcast_channel, NoSocket


def home(request):
    return render(request, 'webhooks/index.html', {})


def new(request):
    webhook_id = str(uuid.uuid4())
    template_args = {'webhook_id': webhook_id}
    return redirect(reverse('webhooks.views.webhook', kwargs=template_args))


@csrf_exempt
def webhook(request, webhook_id):
    if request.GET or request.method != 'GET':
        post = nice_dict(request.POST) if request.POST else scalar_or_list(request.body)
        get = nice_dict(request.GET)
        headers = nice_dict({k: v for k, v in request.META.iteritems() if k.startswith('HTTP_')})
        data = {
            "action": "request:new",
            "method": request.method,
            "ip": request.META.get('REMOTE_ADDR'),
            "host": request.META.get('REMOTE_HOST'),
            "length": request.META.get('CONTENT_LENGTH'),
            "type": request.META.get('CONTENT_TYPE'),
            "post": post,
            "get": get,
            "headers": headers,
        }

        # import json
        # return HttpResponse(json.dumps(data, indent=4, separators=(',', ': ')), content_type='text/plain')

        try:
            broadcast_channel(data, channel=webhook_id)
            return HttpResponse("Message Logged")
        except NoSocket:
            return HttpResponse("Messaged ignored due to no clients watching")

    return render(request, 'webhooks/webhook.html', {'webhook_id': webhook_id, "full_path": request.build_absolute_uri()})
