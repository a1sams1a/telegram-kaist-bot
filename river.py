#!/usr/bin/env python
# -*- coding: utf-8 -*-
from bs4 import BeautifulSoup
import json, requests, datetime

r = requests.get('http://www.koreawqi.go.kr/wQSCHomeLayout_D.wq?action_type=T')
r.encoding = 'euc-kr'

soup = BeautifulSoup(r.text, "lxml")
result = soup.select('#div_layer_btn2_r0 > table > tr')

guri = ''
gapchun = ''
for tr in result:
  title = tr.select('td.start')[0].get_text()
  if (title == u'구리'):
    guri = tr.select('td > span')[0].get_text()
  elif (title == u'갑천'):
    gapchun = tr.select('td > span')[0].get_text()

time = str(datetime.datetime.now().time().hour)

f = open('extdata_river.txt', 'w')
f.write(time + ',' + guri + ',' + gapchun)
f.close()
