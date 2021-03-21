from django.contrib import admin

from forum import models

class BoardAdmin(admin.ModelAdmin):
	list_display = ('title', 'description')

class TopicAdmin(admin.ModelAdmin):
	list_display = ('title', 'board', 'author')

class PostAdmin(admin.ModelAdmin):
	list_display = ('topic', 'author', 'message')

admin.site.register(models.Board, BoardAdmin)
admin.site.register(models.Topic, TopicAdmin)
admin.site.register(models.Post, PostAdmin)
