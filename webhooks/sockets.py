import logging

from socketio.namespace import BaseNamespace
from socketio.mixins import RoomsMixin, BroadcastMixin
from webhooks.sdjango import namespace


@namespace('/webhooks')
class WebhookNamespace(BaseNamespace, RoomsMixin, BroadcastMixin):
    def initialize(self):
        self.logger = logging.getLogger("socketio.webhook")
        self.log("WebhookNamespace socketio session started: %s" % self.socket)

    def log(self, message):
        self.logger.info("[{0}] {1}".format(self.socket.sessid, message))

    def on_join(self, room):
        self.room = room
        self.join(room)
        return True

    def recv_disconnect(self):
        # Remove nickname from the list.
        self.log('Disconnected')
        self.disconnect(silent=True)
        return True


