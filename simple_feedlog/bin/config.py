
class CONFIG:
    debug = True
    max_items = 50
    output_dir = '.'
    database_file = 'news.dat'
    latest_file = 'index.html'
    template = 'template.html'

    feeds = {
        'Google': (
            'http://news.google.com/?output=rss',
        ),
        'Reddit': (
            'http://www.reddit.com/.rss',
        ),
    }
