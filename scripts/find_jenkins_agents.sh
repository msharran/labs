AGENT_NAME="agent2"
FROM_DATE="2023-04-27 05:00:00"
TO_DATE="2023-04-27 05:10:00"

# Helper function to convert date to timestamp
dateToTimestamp() {
  date -u -d "$1" +%s
}
from_timestamp=$(dateToTimestamp "${FROM_DATE}")
to_timestamp=$(dateToTimestamp "${TO_DATE}")

# find all job folders that were created today
  for build_dir in `find jobs/*/builds -mindepth 1 -maxdepth 1 -type d -daystart -ctime 0 -print`; do
    # Check if build.xml exists
    if [ -f "${build_dir}/build.xml" ]; then
      # Extract agent name and timestamp from build.xml
executed_agent=$(cat ${build_dir}/build.xml | grep -e "<node>agent.*</node>" | sed -e 's%<node>%%' -e 's%</node>%%' | awk '{$1=$1};1' | sort | uniq | tr '\n' ',' | sed 's/,$/\n/')
      timestamp_seconds=$(grep -oPm1 "(?<=<timestamp>)[^<]+" "${build_dir}/build.xml")

      # Compare agent name and timestamps
     if echo $executed_agent |  grep -q -w $AGENT_NAME; then
      if [ "${timestamp_seconds:0:-3}" -ge "${from_timestamp}" ] && [ "${timestamp_seconds:0:-3}" -lt "${to_timestamp}" ]; then
              job_name="$build_dir"
        build_number=$(basename "${build_dir}")
        start_time=$(date -u -d @${timestamp_seconds:0:-3} "+%Y-%m-%d %H:%M:%S")
        echo "${job_name} #${build_number} (${start_time})"
      fi
      fi
    fi
  done
