#!/usr/bin/env python
# -*- coding: utf-8 -*-
from bs4 import BeautifulSoup
import json, requests, datetime

r = requests.get('http://www.koreawqi.go.kr/wQSCHomeLayout_D.wq?action_type=T')
r.encoding = 'euc-kr'

soup = BeautifulSoup(r.text, "lxml")

guri = ''
result1 = soup.select('#div_layer_btn2_r1 > table > tr')
for tr in result1:
    title = tr.select('td.start')[0].get_text()
    if (title == u'구리'):
        guri = tr.select('td')[1].get_text()

gapchun = ''
result2 = soup.select('#div_layer_btn2_r3 > table > tr')
for tr in result2:
    title = tr.select('td.start')[0].get_text()
    if (title == u'갑천'):
        gapchun = tr.select('td')[1].get_text()

time = str(datetime.datetime.now().time().hour)

f = open('extdata_river.txt', 'w')
f.write(time + ',' + guri + ',' + gapchun)
f.close()
