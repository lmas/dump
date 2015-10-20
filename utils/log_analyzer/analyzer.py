#!/usr/bin/env python
# Analyze logfiles and show some stats

import sys
import re

class Analyzer:
    def __init__(self):
        #self.re = re.compile(r'(\w{3} \d+ \d+:\d+:\d+) (\w+) (.+?): \[.+?\] (.+)')
        self.re = re.compile(r'(\w{3} \d+ \d+:\d+:\d+) (\w+) (.+?): (.+)')

    def analyze(self, logfile):
        dump = dict()

        total_lines = 0
        with open(logfile, 'r') as f:
            line = f.readline()
            while line:
                total_lines += 1

                tmp = self.re.match(line)
                if tmp:
                    time, host, huh, msg = tmp.groups()
                    dump[msg] = dump.get(msg, 0) + 1

                line = f.readline()

        sort = sorted(dump.items(), key=lambda x: x[1])
        sort.reverse()
        return total_lines, sort

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print('ERROR: Must provide path to logfile!')
        sys.exit(1)

    logfile = sys.argv[1]
    print('Analyzing %s...' % logfile)

    test = Analyzer()
    totals, dump = test.analyze(logfile)
    unique = len(dump)

    print('Lines analyzed:', totals)
    print('Unique lines:', unique)
    print('\nPercent\tHits\tMessage')

    for line, hits in dump:
        percent = hits / float(unique)
        print('%.3f\t%i\t%s' %  (percent, hits, line))
