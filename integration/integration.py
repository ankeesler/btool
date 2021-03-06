import subprocess
import sys
import unittest

registry = "https://btoolregistry.cfapps.io"
btool_in = "btool"
btool_out = "/tmp/btool"

class TestBtool(unittest.TestCase):
    def test_list(self):
        ac_output = subprocess.check_output([btool_in, "-root", "source", "-registry", registry, "-target", "btool", "-list"])
        self.assertTrue(". source/btool.o" in ac_output)
        self.assertTrue(". source/app/app.o" in ac_output)

    def test_build_btool(self):
        subprocess.check_call([btool_in, "-root", "source", "-registry", registry, "-target", "btool"])
        subprocess.check_call(["mv", "source/btool", btool_out])
        subprocess.check_call([btool_in, "-root", "source", "-registry", registry, "-target", "btool", "-clean"])

    def test_build_example_basic_c(self):
        subprocess.check_call([btool_in, "-root", "example/BasicC", "-registry", registry, "-target", "main"])
        subprocess.check_call(["example/BasicC/main"])
        subprocess.check_call([btool_in, "-root", "example/BasicC", "-registry", registry, "-target", "main", "-clean"])

    def test_build_example_basic_cc(self):
        subprocess.check_call([btool_in, "-root", "example/BasicCC", "-registry", registry, "-target", "main"])
        subprocess.check_call(["example/BasicCC/main"])
        subprocess.check_call([btool_in, "-root", "example/BasicCC", "-registry", registry, "-target", "main", "-clean"])

    def test_test(self):
        subprocess.check_call([btool_in, "-registry", registry, "-root", "example/BasicCC", "-target", "dep-1/dep-1-test", "-run"])
        subprocess.check_call(["example/BasicCC/dep-1/dep-1-test"])
        subprocess.check_call([btool_in, "-registry", registry, "-root", "example/BasicCC", "-target", "dep-1/dep-1-test", "-clean"])

if __name__ == '__main__':
    if len(sys.argv) > 1:
        btool_in = sys.argv[1]
    if len(sys.argv) > 2:
        btool_out = sys.argv[2]
    if len(sys.argv) > 3:
        registry = sys.argv[3]

    suite = unittest.TestLoader().loadTestsFromTestCase(TestBtool)
    unittest.TextTestRunner().run(suite)
