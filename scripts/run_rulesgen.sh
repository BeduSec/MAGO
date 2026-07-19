# Copyright (c) BeduSec. All rights reserved.
#!/bin/bash
cd tools/rulesgen
python cli.py generate --policy ../../examples/policy.yaml --out ../../examples/rulesets/default.json
python cli.py test --rules ../../examples/rulesets/default.json --target http://localhost:8080