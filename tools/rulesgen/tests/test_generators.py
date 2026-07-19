# Copyright (c) BeduSec. All rights reserved.
import json
import os
import tempfile
import requests
from generators import generate_ruleset

def run_tests(rules_path, target):
    with open(rules_path) as f:
        rules = json.load(f)
    print(f"Loaded {len(rules)} rules. Running basic validation...")
    blocked = False
    try:
        resp = requests.get(target + '/', timeout=2)
        if resp.status_code == 403:
            blocked = True
            print("Request blocked (WAF) - OK")
    except requests.ConnectionError:
        print("MAGO not reachable, skipping live test")
    print("Tests completed.")

def test_generate():
    policy_yaml = """
rules:
  - id: test-rule
    priority: 1
    action: block
    conditions:
      - field: path
        operator: regex
        value: "/test"
    match_type: all
"""
    with tempfile.NamedTemporaryFile(mode='w', suffix='.yaml', delete=False) as pf:
        pf.write(policy_yaml)
        policy_path = pf.name
    out_path = tempfile.mktemp(suffix='.json')
    try:
        generate_ruleset(policy_path, out_path)
        with open(out_path) as f:
            rules = json.load(f)
        assert len(rules) == 1
        assert rules[0]['id'] == 'test-rule'
    finally:
        os.unlink(policy_path)
        os.unlink(out_path)