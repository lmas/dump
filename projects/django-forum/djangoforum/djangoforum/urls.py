from django.conf.urls import patterns, include, url

from django.contrib import admin
admin.autodiscover()

from django.shortcuts import redirect

urlpatterns = patterns('',
    # Examples:
    # url(r'^$', 'djangoforum.views.home', name='home'),
    # url(r'^blog/', include('blog.urls')),

    url(r'^admin/', include(admin.site.urls)),

    url(r'^forum/', include('forum.urls')),
)
