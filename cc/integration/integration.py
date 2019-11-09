import subprocess
import sys
import unittest

btool = "btool"

class TestBtool(unittest.TestCase):
    def test_list(self):
        subprocess.check_call([btool, "-target", "btool"])

        ac_output = subprocess.check_output(["./btool", "-target", "btool", "-list"])
        self.assertTrue(". ./btool.o" in ac_output)
        self.assertTrue(". ./app/app.o" in ac_output)

    def test_build(self):
        subprocess.check_call([btool, "-target", "btool"])
        subprocess.check_call(["mv", "btool", "/tmp/btool-from-go"])
        subprocess.check_call(["./script/clean.sh"])

        subprocess.check_call(["/tmp/btool-from-go", "-target", "btool"])
        subprocess.check_call(["mv", "btool", "/tmp/btool-from-cc"])
        subprocess.check_call(["./script/clean.sh"])

        subprocess.check_call(["/tmp/btool-from-cc", "-target", "btool"])
        subprocess.check_call(["mv", "btool", "/tmp/btool-from-cc-from-cc"])
        subprocess.check_call(["./script/clean.sh"])

        subprocess.check_call(["/tmp/btool-from-cc-from-cc", "-target", "btool"])
        subprocess.check_call(["ls", "btool"])

if __name__ == '__main__':
    if len(sys.argv) > 1:
        btool = sys.argv[1]
    suite = unittest.TestLoader().loadTestsFromTestCase(TestBtool)
    unittest.TextTestRunner().run(suite)
