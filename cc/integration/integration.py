import subprocess
import unittest

class TestBtool(unittest.TestCase):
    def test_list(self):
        subprocess.check_call(["btool", "-target", "btool"])

        ac_output = subprocess.check_output(["./btool", "-target", "btool", "-list"])
        self.assertTrue(". ./btool.o" in ac_output)
        self.assertTrue(". ./app/app.o" in ac_output)

    def test_build(self):
        subprocess.check_call(["btool", "-target", "btool"])
        subprocess.check_call(["mv", "btool", "/tmp/btool-from-go"])
        subprocess.check_call(["./script/clean.sh"])

        subprocess.check_call(["/tmp/btool-from-go", "-target", "btool"])
        subprocess.check_call(["mv", "btool", "/tmp/btool-from-cc"])
        subprocess.check_call(["./script/clean.sh"])

        subprocess.check_call(["/tmp/btool-from-cc", "-target", "btool"])

        cs0 = subprocess.check_output(["md5", "-q", "/tmp/btool-from-cc"])
        cs1 = subprocess.check_output(["md5", "-q", "btool"])
        self.assertEqual(cs0, cs1)

if __name__ == '__main__':
    unittest.main()
