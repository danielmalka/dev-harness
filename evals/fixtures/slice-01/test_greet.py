import unittest

from greet import greet


class TestGreet(unittest.TestCase):
    def test_says_hello(self):
        self.assertEqual(greet("Ada"), "Hello, Ada")

    def test_rejects_blank(self):
        with self.assertRaises(ValueError):
            greet("  ")

    def test_rejects_empty(self):
        with self.assertRaises(ValueError):
            greet("")


if __name__ == "__main__":
    unittest.main()
