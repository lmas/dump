from django.conf.urls import patterns, url

from forum import views

urlpatterns = patterns('',
    # Examples:
    # url(r'^$', 'djangoforum.views.home', name='forum_index'),

	url(r'^$', views.BoardList.as_view(), name='forum_board_list'),

	url(r'^board/(?P<board>\d+)$', views.TopicList.as_view(), name='forum_topic_list'),
	url(r'^board/(?P<board>\d+)/new$', views.TopicCreate.as_view(), name='forum_topic_create'),

	url(r'^topic/(?P<topic>\d+)$', views.PostList.as_view(), name='forum_post_list'),
)
