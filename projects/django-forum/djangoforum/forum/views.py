from django.shortcuts import render, get_object_or_404
from django.views.generic import ListView, DetailView
from django.views.generic import CreateView

from forum.models import Board, Topic, Post

class BoardList(ListView):
        model = Board

class GetBoardMixin(object):
        def get_context_data(self, **kwargs):
                self.board = get_object_or_404(Board, pk=self.kwargs['board'])
                context = super(GetBoardMixin, self).get_context_data(**kwargs)
                context['board'] = self.board
                return context

class TopicList(GetBoardMixin, ListView):
        def get_queryset(self):
                self.board = get_object_or_404(Board, pk=self.kwargs['board'])
                return Topic.objects.filter(board=self.board)

class TopicCreate(GetBoardMixin, CreateView):
        model = Topic
        fields = ('title',)

class PostList(ListView):
        def get_queryset(self):
                self.topic = get_object_or_404(Topic, pk=self.kwargs['topic'])
                return Post.objects.filter(topic=self.topic)

        def get_context_data(self, **kwargs):
                context = super(PostList, self).get_context_data(**kwargs)
                context['topic'] = self.topic
                return context

