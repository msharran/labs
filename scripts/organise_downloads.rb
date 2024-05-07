require 'fileutils'

downloads = Dir.glob("#{Dir.home}/Downloads/202{3,4}-*")

# # group by date, then by file type
# # directory structure: ~/Downloads/2014-01-01/zip
# # 2014-01-01/jpg
#
# downloads.each do |d|
#   next if File.directory? d
#   mtime = File.mtime d
#   date = mtime.strftime "%Y-%m-%d"
#   type = d.split('.').last
#   dir = "#{Dir.home}/Downloads/#{type}"
#   FileUtils.mkdir_p dir
#   FileUtils.mv d, dir
# end
  
