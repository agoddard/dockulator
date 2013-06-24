#!/usr/bin/env ruby
left, op, right = ARGV[0].split

allowed = ["+","-","*","/"]

if allowed.include? op
  answer = left.to_f.send(op, right.to_f)
  puts answer
else
  exit 1
  puts "Error"
end
