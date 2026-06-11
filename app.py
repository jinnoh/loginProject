import os

from flask import Flask, render_template_string, request, redirect, url_for, session

app = Flask(__name__)
