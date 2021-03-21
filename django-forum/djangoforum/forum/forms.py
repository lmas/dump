
from django import forms

from forum.models import Board, Topic, Post

class TopicCreateForm(forms.ModelForm):
    class meta:
        model = Topic
        fields = ('title',)
