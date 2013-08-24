from django.conf.urls import patterns, include, url

# Uncomment the next two lines to enable the admin:
# from django.contrib import admin
# admin.autodiscover()

urlpatterns = patterns('',
    # Examples:
    # url(r'^$', 'webhookr.views.home', name='home'),
    # url(r'^webhookr/', include('webhookr.foo.urls')),

    # Uncomment the admin/doc line below to enable admin documentation:
    # url(r'^admin/doc/', include('django.contrib.admindocs.urls')),

    # Uncomment the next line to enable the admin:
    # url(r'^admin/', include(admin.site.urls)),
    url(r'^$', 'webhooks.views.home', name='home'),
    url(r'^(?P<webhook_id>[\S-]{36})/?.*$', 'webhooks.views.webhook', name='webhook'),
    url(r'^new/?$', 'webhooks.views.new', name='new'),
    url("", include('django_socketio.urls')),
)
