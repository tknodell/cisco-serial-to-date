# This script will take an individual or text file of cisco serial numbers and return the date manufacturered

# The serial number will be in the format: ‘LLLYYWWSSSS’; where ‘YY’ is the year of manufacture and ‘WW’ is the week of manufacture.
# The date code can be found in the 4 middle digits of the serial number
#
# Here are the Manufacturing Year Codes:
#                01 = 1997                              06 = 2002                              11 = 2007                              16 = 2012
#                02 = 1998                              07 = 2003                              12 = 2008                              17 = 2013
#                03 = 1999                              08 = 2004                              13 = 2009                              18 = 2014
#                04 = 2000                              09 = 2005                              14 = 2010                              19 = 2015
#                05 = 2001                              10 = 2006                              15 = 2011                              20 = 2016
#
# And here are the Manufacturing Week Codes:
#                1-5 : January                 15-18 : April                       28-31 : July              41-44 : October
#                6-9 : February                19-22 : May                         32-35 : August            45-48 : November
#                10-14 : March                 23-27 : June                        36-40 : September         49-52 : December

require 'Date'

@years = {
  '01' => 1997, '02' => 1998, '03' => 1999, '04' => 2000, '05' => 2001, '06' => 2002, '07' => 2003, '08' => 2004, '09' => 2005, '10' => 2006,
  '11' => 2007, '12' => 2008, '13' => 2009, '14' => 2010, '15' => 2011, '16' => 2012, '17' => 2013, '18' => 2014, '19' => 2015, '20' => 2016
}

def serialtodate(num)
  year = @years[num[3..4]]
  month = Date.commercial(year, num[5..6].to_i).strftime('%m')
  # Assume Monday of given week number
  week_start = Date.commercial(year, num[5..6].to_i, 1).strftime('%d')
  t = Time.new(year, month, week_start)
  puts num
  puts t.strftime('%Y-%m-%d')
end

unless ARGV.length == 1
  puts "You must specify either a single serial number, or a file with a list of serial numbers to process\n\n"
  puts "Example: ruby serial_to_date.rb FAA04459FNI\n"
  puts "Example: ruby serial_to_date.rb cisco_serials.txt\n"
  exit
end

input = ARGV.first

if File.exist?(input)
  File.open(input, 'r') do |f|
    f.each_line do |line|
      puts serialtodate(line) unless line.chomp.empty?
    end
  end
else
  puts serialtodate(input.to_s)
end
