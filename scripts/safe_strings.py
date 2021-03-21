
alpha_nums = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ'

def clean_string(data):
	return ''.join(c for c in data if c in valid_chars)

#import string
#trans_table = string.maketrans('','')
#valid_chars = '._-%s' % alpha_nums
#def valid_string(s):
#	return not s.translate(trans_table, valid_chars)

def html_encode(data):
	data = data.replace('&', '&amp;')
	data = data.replace('<', '&lt;')
	data = data.replace('>', '&gt;')
	data = data.replace('"', '&quot;')
	data = data.replace('\'', '&#39;')
	data = data.replace('`', '&#96;')
	data = bytes(data, 'utf-8')
	return data

