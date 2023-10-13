require 'net/http'
require 'json'
require 'yaml'

def load_plugins(file_name)
  plugins_file = YAML.load_file(file_name)
  plugin_names = []

  plugins_file['plugins'].each do |plugin|
    name = plugin['name']
    plugin_names.push(name)
  end

  plugin_names
end

def get_latest_plugin_version(plugin_name)
  url = "https://plugins.jenkins.io/api/plugin/#{plugin_name}"
  uri = URI(url)
  response = Net::HTTP.get(uri)
  json_data = JSON.parse(response)
  json_data['version']
rescue StandardError => e
  puts "Error: #{e.message}"
  nil
end

def save_plugins(file_name, plugins)
  File.open(file_name, mode: 'w') do |f|
    f.write(plugins.to_yaml)
  end
end

# MAIN

puts 'plugins.yaml file path is required as argument' if ARGV.count.zero?

file = ARGV[0]
plugin_names = load_plugins(file)

puts "=> loaded #{plugin_names.count} plugins from #{file}"

plugins_with_latest_versions = {}

plugin_names.each do |plugin_name|
  version = get_latest_plugin_version(plugin_name)
  plugins_with_latest_versions[plugin_name] = version if version
  puts "==> downloaded latest version for #{plugin_name}"
end

new_versions = []
puts '=> latest Jenkins plugin versions:'
plugins_with_latest_versions.each do |plugin_name, version|
  new_versions.push({ 'name' => plugin_name, 'version' => version })
end

puts '=> saving to plugins.yml'
save_plugins('plugins.yml', { 'plugins' => new_versions })
