import logging

from socketio.namespace import BaseNamespace
from socketio.mixins import RoomsMixin, BroadcastMixin


class WebhookrChannelMixin(object):
    room_key = 'rooms'

    def __init__(self, *args, **kwargs):
        super(WebhookrChannelMixin, self).__init__(*args, **kwargs)
        if self.room_key not in self.session:
            self.session[self.room_key] = set()  # a set of simple strings

    def join(self, room):
        """Lets a user join a room on a specific Namespace."""
        self.session[self.room_key].add(self._get_room_name(room))

    def leave(self, room):
        """Lets a user leave a room on a specific Namespace."""
        self.session[self.room_key].remove(self._get_room_name(room))

    def _get_room_name(self, room):
        return self.ns_name + '_' + room

    def room_subscribers(self, room, include_self=False):
        room_name = self._get_room_name(room)

        for sessid, socket in self.socket.server.sockets.iteritems():
            if self.room_key not in socket.session:
                continue
            if room_name in socket.session[self.room_key] and (include_self or self.socket != socket):
                yield (sessid, socket)

    def all_rooms(self):
        return (x[len(self.ns_name):] for x in self.session.get(self.room_key, []))

    def _emit_to_channel(self, room, event, include_self, *args):
        """This is sent to all in the room (in this particular Namespace)"""
        message = dict(type="event", name=event, args=args, endpoint=self.ns_name)

        for sessid, socket in self.room_subscribers(room, include_self=include_self):
            socket.send_packet(message)

    def emit_to_channel(self, room, event, *args):
        self._emit_to_channel(room, event, False, *args)

    def emit_to_channel_and_me(self, room, event, *args):
        self._emit_to_channel(room, event, True, *args)


class WebhookNamespace(BaseNamespace, WebhookrChannelMixin, BroadcastMixin):
    def initialize(self):
        self.logger = logging.getLogger("socketio.webhook")
        self.log("WebhookNamespace socketio session started: %s" % self.socket)

    def log(self, message):
        self.logger.info("[{0}] {1}".format(self.socket.sessid, message))

    def emit_subscriber_count(self, room):
        # Enumerate to get length of subscribers while being lazy
        i = 0
        for i, x in enumerate(self.room_subscribers(room, include_self=True)):
            self.logger.debug('[emit_subscriber_count] i= {}'.format(i))
        total_subscribers = i + 1
        self.log('Emitting totalSubscribers for {}: {}'.format(room, total_subscribers))
        self.emit_to_channel_and_me(room, 'subscriber_joined', {'totalSubscribers': total_subscribers})

    def on_join(self, room):
        self.log('Connected')
        self.room = room
        self.join(room)
        self.emit_subscriber_count(room)

        return True

    def recv_disconnect(self):
        # Remove nickname from the list.
        self.log('Disconnected')
        for room in self.all_rooms():
            self.emit_subscriber_count(room)
        self.disconnect(silent=True)
        return True

