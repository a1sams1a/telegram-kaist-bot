#!/usr/bin/env python
# -*- coding: utf-8 -*-
from bs4 import BeautifulSoup
import requests

r = requests.get('http://www.kma.go.kr/weather/forecast/timeseries.jsp?searchType=SETINFO&setinfoCode=3020054000')
r.encoding = 'euc-kr'

soup = BeautifulSoup(r.text, "lxml")
title = soup.select('td.Situbg > dl > dt.bold.Wimg')[0].get_text().strip()
temp = soup.select('td.Situbg > dl > dd.hum')[0].get_text().strip()
low_temp = soup.select('table.forecastNew3 > tbody > tr > td.bg_tomorrow > .low_deg')[0].get_text().strip()
high_temp = soup.select('table.forecastNew3 > tbody > tr > td.bg_tomorrow > .high_deg')[0].get_text().strip()

f = open('extdata_weather.txt', 'w')
f.write((title + ',' + temp + ',' + high_temp + ',' + low_temp).encode('utf-8'))
f.close()
