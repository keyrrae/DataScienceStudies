from __future__ import absolute_import
import sys
sys.path.insert(0, "/Users/xuanwang/Applications/buck/third-party/py/pathlib")
sys.path.insert(0, "/Users/xuanwang/Applications/buck/third-party/py/pywatchman")
sys.path.insert(0, "/Users/xuanwang/Applications/buck/third-party/py/typing/python2")
sys.path.insert(0, "/Users/xuanwang/Applications/buck/build/classes")
sys.path.insert(0, "/Users/xuanwang/Workspace/goprojects/src/github.com/keyrrae/monimenta_backend/.buckd/tmp/buck_run.kdjcxz/buck_python_program2026661400832139245")
if __name__ == '__main__':
    try:
        from buck_parser import buck
        buck.main()
    except KeyboardInterrupt:
        print >> sys.stderr, 'Killed by User'
