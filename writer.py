import csv
import random
import string

def projNum_generator():
  return ''.join(random.choice(string.digits) for _ in range(10))

def user_generator():
  size = 5
  userchars = string.ascii_lowercase + string.digits
  return ''.join(random.choice(userchars) for _ in range(size))

def week_generator():
  return ''.join(random.choice(string.digits) for _ in range(2))

def time_generator():
  return str(round(random.uniform(0, 10), 1))

def createRow():
  week = week_generator()
  user = user_generator()
  mon = time_generator()
  tue = time_generator()
  wed = time_generator()
  thu = time_generator()
  fri = time_generator()
  projNum = projNum_generator()
  
  row = (week, user, mon, tue, wed, thu, fri, projNum)
  return row

with open('text.csv', 'w', newline='') as csvfile:
    writer = csv.writer(csvfile, dialect='excel')
    for i in range(1000000):
      writer.writerow(createRow())


# def projNum_list(n):
#   my_projNums = []
#   for i in range(n):
#     my_projNums.append(projNum_generator())
#   return my_projNums

# def user_list(n):
#   my_users = []
#   for i in range(n):
#     my_users.append(user_generator())
#   return my_users

# def week_list():
#   weeklist = []
#   for i in range(1, 53):
#     weeklist.append(i)
#   return weeklist