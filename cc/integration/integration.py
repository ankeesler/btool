import subprocess
import sys
import unittest

registry = "http://btoolregistry.cfapps.io"
btool_in = "btool"
btool_out = "/tmp/btool"

class TestBtool(unittest.TestCase):
    def test_list(self):
        ac_output = subprocess.check_output([btool_in, "-registry", registry, "-target", "btool", "-list", "-loglevel", "debug"])
        self.assertTrue(". btool.o" in ac_output or ". ./btool.o" in ac_output)
        self.assertTrue(". app/app.o" in ac_output or ". ./app/app.o" in ac_output)

    def test_build(self):
        subprocess.check_call([btool_in, "-registry", registry, "-target", "btool", "-loglevel", "debug"])
        subprocess.check_call(["mv", "btool", btool_out])
        subprocess.check_call([btool_in, "-registry", registry, "-target", "btool", "-loglevel", "debug", "-clean"])

    def test_test(self):
        subprocess.check_call([btool_in, "-registry", registry, "-root", "../example/BasicCC", "-target", "dep-1/dep-1-test", "-loglevel", "debug", "-run"])
        subprocess.check_call(["../example/BasicCC/dep-1/dep-1-test"])
        subprocess.check_call([btool_in, "-registry", registry, "-root", "../example/BasicCC", "-target", "dep-1/dep-1-test", "-loglevel", "debug", "-clean"])

if __name__ == '__main__':
    if len(sys.argv) > 1:
        btool_in = sys.argv[1]
    if len(sys.argv) > 2:
        btool_out = sys.argv[2]
    if len(sys.argv) > 3:
        registry = sys.argv[3]

    suite = unittest.TestLoader().loadTestsFromTestCase(TestBtool)
    unittest.TextTestRunner().run(suite)
