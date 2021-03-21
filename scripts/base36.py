# from: http://webpy.org/

def base36(q):
	"""
	Converts an integer to base 36 (a useful scheme for human-sayable IDs).

		>>> to36(35)
		'z'
		>>> to36(119292)
		'2k1o'
		>>> int(to36(939387374), 36)
		939387374
		>>> to36(0)
		'0'
		>>> to36(-393)
		Traceback (most recent call last):
			...
		ValueError: must supply a positive integer

	"""
	if q < 0: raise ValueError, "must supply a positive integer"
	letters = "0123456789abcdefghijklmnopqrstuvwxyz"
	converted = []
	while q != 0:
		q, r = divmod(q, 36)
		converted.insert(0, letters[r])
	return "".join(converted) or '0'
