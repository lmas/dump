from django.db import models
from django.contrib.auth.models import User

# TODO:
# timestamps

class Board(models.Model):
	title = models.CharField(max_length=255, unique=True)
	description = models.TextField()

	def __unicode__(self):
		return self.title

class Topic(models.Model):
	author = models.ForeignKey(User)
	board = models.ForeignKey(Board)
	title = models.CharField(max_length=255)
	created = models.DateTimeField(auto_now_add=True)

	def __unicode__(self):
		return self.title

class Post(models.Model):
	author = models.ForeignKey(User)
	topic = models.ForeignKey(Topic)
	message = models.TextField()
	created = models.DateTimeField(auto_now_add=True)
	modified = models.DateTimeField(auto_now=True)

	def __unicode__(self):
		return '%i:%s' % (self.pk, self.message)

