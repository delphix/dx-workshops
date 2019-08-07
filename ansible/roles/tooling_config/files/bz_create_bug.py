#! /usr/bin/env python
#Adam Bowen - Oct 2017
#This script creates a bug in Bugzilla
# --- Create vFiles VDB
#requirements
#pip install docopt requests

#The below doc follows the POSIX compliant standards and allows us to use 
#this doc to also define our arguments for the script.


"""Create a bug in bugzilla
Usage:
  bz_create_bug.py --hostname <name_or_ip> (--api_key <key> | --login <username> --password <password>) --summary <summary> --description <description>
  bz_create_bug.py -h | --help | -v | --version
Provision VDB from a defined source on the defined target environment.

Examples:
  bz_create_bug.py --hostname 52.15.144.164 --api_key AYpCL5uSbGFhQMUq6sSWcYo2qNoQskIxdL0EcK0L --summary 'Selenium Test Results for Project X' --description 'Failed test results go here'
  bz_create_bug.py --hostname 18.216.206.77 --login admin --password --summary 'Selenium Test Results for Project X' --description 'Failed test results go here'

Options:
  --hostname <name_or_ip>     The IP of the bugzilla server.
  --api_key <key>             The API key used to authenticate
  --login <username>          The bugzilla username to use
  --password <password>       The bugzilla password to use
  --summary <summary>         The summary for the bug
  --description <description> The detailed description of the bug
  -h --help                   Show this screen.
  -v --version                Show version.
"""
 
VERSION = 'v0.02'


import requests
import json
from os.path import basename
from docopt import docopt

def main(argv):
  bz_url = 'http://'+ arguments['--hostname'] + '/bugzilla/rest/bug'
  bz_api_key = arguments['--api_key']
  summary = arguments['--summary']
  description = arguments['--description']
  bz_username = arguments['--login']
  bz_password = arguments['--password']

  if bz_api_key != None:
    bz_params = {'api_key': bz_api_key}
  else:
    bz_params = {'login': bz_username,'password': bz_password}

 # r = requests.post(bz_url, params={'api_key': bz_api_key}, json={
  r = requests.post(bz_url, params=bz_params, json={
                  "product" : "TestProduct",
                  "component" : "TestComponent",
                  "version" : "unspecified",
                  "summary" : summary,
                  "op_sys" : "All",
                  "platform" : "All",
                  "severity" : "critical",
                  "description" : description
                   })
  if r.ok:
    for k, v in r.json().items():
      print "BUG " + str(v)
  else:
    print "Something went wrong"
    print r.status_code
    print r.reason
    print r.content
if __name__ == "__main__":
    #Grab our arguments from the doc at the top of the script
    arguments = docopt(__doc__, version=basename(__file__) + " " + VERSION)
    #Feed our arguments to the main function, and off we go!
    main(arguments)