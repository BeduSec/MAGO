# Copyright (c) BeduSec. All rights reserved.
import argparse
import json
import sys
from generators import generate_ruleset
from tests.test_generators import run_tests

def main():
    parser = argparse.ArgumentParser(prog='rulesgen', description='MAGO ruleset generator and test harness')
    subparsers = parser.add_subparsers(dest='command', required=True)

    gen_parser = subparsers.add_parser('generate', help='Generate ruleset from policy YAML')
    gen_parser.add_argument('--policy', required=True, help='Path to policy YAML file')
    gen_parser.add_argument('--out', required=True, help='Output JSON ruleset file')

    test_parser = subparsers.add_parser('test', help='Test ruleset against running MAGO instance')
    test_parser.add_argument('--rules', required=True, help='Path to ruleset JSON')
    test_parser.add_argument('--target', default='http://localhost:8080', help='MAGO base URL')

    args = parser.parse_args()

    if args.command == 'generate':
        generate_ruleset(args.policy, args.out)
        print(f"Ruleset written to {args.out}")
    elif args.command == 'test':
        run_tests(args.rules, args.target)

if __name__ == '__main__':
    main()