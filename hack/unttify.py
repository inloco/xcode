#!/usr/bin/env python3
import fcntl
import os
import pty
import re
import struct
import subprocess
import sys
import threading
import tty


def io(master_fd):
    fcntl.ioctl(master_fd, tty.TIOCSWINSZ, struct.pack('HHHH', 24, 80, 0, 0))

    with os.fdopen(master_fd) as stdout:
        cs = re.compile(r'\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])')
        for l in iter(stdout.readline, ''):
            sys.stdout.write(cs.sub('', l))
            sys.stdout.flush()


def main():
    master_fd, slave_fd = pty.openpty()

    t = threading.Thread(target=io, args=(master_fd,))
    t.setDaemon(True)
    t.start()

    proc = subprocess.Popen(sys.argv[1:], stdin=slave_fd, stdout=slave_fd, stderr=slave_fd)
    sys.exit(proc.wait())


if __name__ == '__main__':
    main()
