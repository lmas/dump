#!/usr/bin/env python

from hashlib import sha256 as hasher

# The target string each block must begin with to be valid
TARGET_BYTES_LENGTH = 5
TARGET_BYTES = '0' * TARGET_BYTES_LENGTH

# Hard limit before we give up finding a new block
MAX_ITERATIONS = 10000000

def find_block(data):
    '''Try to find a block (based on data) containing the target bytes.'''
    for nonce in range(MAX_ITERATIONS):
        # create a payload we can create a new block from
        payload = ''.join((str(data), str(nonce))).encode('utf-8')
        block = hasher(payload).hexdigest()

        # see if the new block is a valid one
        if block[0:TARGET_BYTES_LENGTH] == TARGET_BYTES:
            return block

class BlockChain(object):
    def __init__(self, data):
        '''Initiate the chain with some random data from caller.'''
        self.initial_data = data
        self.blocks = [data]
        self.__end = False

    def count(self):
        '''Return the amount of found blocks in this chain.'''
        return len(self.blocks)-1

    def generate(self):
        '''Generate a new block for the chain.'''
        # Prevent running if it's the end of the chain
        if self.__end:
            return

        new_block = find_block(self.blocks)
        if new_block == None:
            # No new blocks found, it's the end of this chain
            self.__end = True
            return

        self.blocks.append(new_block)
        return new_block

###############################################################################

def main():
    chain = BlockChain(1)
    try:
        while True:
            print('\rBlocks found: {}'.format(chain.count()), end='')
            latest_block = chain.generate()
            if latest_block == None:
                print('\nEnd of chain found.', end='')
                raise KeyboardInterrupt
    except KeyboardInterrupt:
        print('\nBlock chain:')
        for tmp in chain.blocks:
            print(tmp)

if __name__ == '__main__':
    main()

