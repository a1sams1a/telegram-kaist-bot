#!/usr/bin/env python
# -*- coding: utf-8 -*-
from bs4 import BeautifulSoup
import requests

rlist = { 'n':'fclt', 'w':'west', 'e':'east1' }
tlist = [ '아침', '점심', '저녁' ]

for loc in rlist:
    r = requests.get('http://www.kaist.ac.kr/_prog/fodlst/index.php?site_dvs_cd=kr&menu_dvs_cd=050303&dvs_cd=fclt&dvs_cd=' + rlist[loc])
    r.encoding = 'utf-8'

    soup = BeautifulSoup(r.text, "lxml")
    result = soup.select('table.menuTb tbody td')

    f = open('extdata_food_' + loc + '.txt', 'w')
    for i in range(3):
        f.write('===' + tlist[i] + '===\n' + result[i].get_text().strip().encode('utf-8'))
    f.close()
