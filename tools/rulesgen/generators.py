# Copyright (c) BeduSec. All rights reserved.
import json
import yaml

def generate_ruleset(policy_path, output_path):
    with open(policy_path, 'r') as f:
        policy = yaml.safe_load(f)
    rules = policy.get('rules', [])
    with open(output_path, 'w') as out:
        json.dump(rules, out, indent=2)