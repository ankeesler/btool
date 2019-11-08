import subprocess
import unittest

class TestBtool(unittest.TestCase):
    def test_list(self):
        subprocess.check_call(["btool", "-target", "btool"])

        ac_output = subprocess.check_output(["./btool", "-target", "btool", "-list"])
        self.assertTrue(". ./btool.o" in ac_output)
        self.assertTrue(". ./app/app.o" in ac_output)

if __name__ == '__main__':
    unittest.main()
