require 'yaml'

plugins_file = YAML.load_file('./existing-plugins.yaml')

plugins_file['plugins'].each do |plugin|
  puts plugin['name']
end
