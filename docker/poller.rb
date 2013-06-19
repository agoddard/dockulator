require 'mongo/driver'

os = %w(jdoi34j djwoief fj93jg4 jfio23jf)
language = %w(rb pl py sh)


while true do
  puts "polling mongo"
  mongo_uri = 'mongodb://dkkr:joi82h32fhoih3f@dharma.mongohq.com:10075/app16386574'
  calculations = Mongo::Connection.from_uri(mongo_uri).db('app16386574').collection('calculations')

  calculations.all answer: nil do |calculation|
    puts "found #{calculation['_id']}"
    puts "processing #{calculation['calculation']} with #{os.sample} and #{language.sample}"
    puts %x(./docker -d /opt/calculate.#{language} --calculation='#{calculation['calculation']}' --id=#{calculation['_id']})
    #parse response from docker
    #update MongoDB
  end
  sleep 10
end
