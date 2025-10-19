#!/usr/bin/bash
##本人的，你要是运行不了自己改
#
###有错立即退出
##set -euo pipefail
# #不直接退出是因为要分裂清理旧文件和备份过程，如果备份了但是没清理不要直接退出
#
##配置
#BACKUP_DIR=$(cd "$(dirname "$0")/../vote/backup" && pwd)
#MYSQL_USERNAME=$(grep 'username' ../vote/config/database.go | awk -F '"' '{print $2}')
#MYSQL_PASSWORD=$(grep 'password' ../vote/config/database.go | awk -F '"' '{print $2}')
#MYSQL_HOST=$(grep 'host' ../vote/config/database.go | awk -F '"' '{print $2}')
#MYSQL_PORT=$(grep 'port' ../vote/config/database.go | awk -F '"' '{print $2}')
#DB_NAME=$(grep 'Dbname' ../vote/config/database.go | awk -F '"' '{print $2}')
#REMAIN_DAY=3 #备份三天
#
##不存在就建一个
#mkdir -p "$BACKUP_DIR"
#
##带上时间戳方便查
#BACKUP_FILE="$BACKUP_DIR/${DB_NAME}_$(date +%Y%m%d_%H%M%S).sql.gz"
#
#echo "start to backup $DB_NAME"
##命令行参数
#mysqldump -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USERNAME" -p"$MYSQL_PASSWORD" --database "$DB_NAME" | gzip > "$BACKUP_FILE"
##-f检查文件是否存在，du -b可以以字节为单位显示文件大小，cut -f1能提取第一列
#
#if [ -s "$BACKUP_FILE" ]; then
#  echo "success to backup"
#else
#  echo "fail to backup"
#  exit 1
#fi
#
#echo "clear up out-date data"
#find "$BACKUP_DIR" -name "${DB_NAME}_*.sql.gz" -type f -mtime +"$REMAIN_DAY" -delete


//感觉之前那个不是很安全，重写一个吧
set -euo pipefail

#默认的配置
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
DEFAULT_CONFIG_FILE="${SCRIPT_DIR}/backup.conf"
DEFAULT_BACKUP_DIR="${SCRIPT_DIR}/../vote/backup"
DEFAULT_REMAIN_DAY=3

#配置变量
CONFIG_FILE=""
BACKUP_DIR=""
MYSQL_HOST=""
MYSQL_PORT=""
MYSQL_USER=""
MYSQL_PASSWORD=""
DB_NAME=""
REMAIN_DAY=""

parse(){
  while [[ $# -gt 0 ]];do
    case $1 in
      -c|--config)
        CONFIG_FILE="$2"
        shift 2
        ;;
      -e|--env)
        load_from_env
        shift
        ;;
      *)
        echo "unknown command $1" >&2
        exit 1
        ;;
    esac
  done
}

#从环境变量加载
load_from_env(){

  BACKUP_DIR="${BACKUP_DIR:-$DEFAULT_BACKUP_DIR}"
  MYSQL_HOST="${MYSQL_HOST:-localhost}"
  MYSQL_PORT="${MYSQL_PORT:-3306}"
  MYSQL_USER="${MYSQL_USER:-}"
  MYSQL_PASSWORD="${MYSQL_PASSWORD:-}"
  DB_NAME="${DB_NAME:-}"
  REMAIN_DAY="${REMAIN_DAY:-$DEFAULT_REMAIN_DAY}"
}

# 从配置文件加载配置
load_from_config() {
    local config_file="$1"

    if [[ ! -f "$config_file" ]]; then
        echo "fail to find file: $config_file" >&2
        return 1
    fi


    while IFS='=' read -r key value || [[ -n "$key" ]]; do
        # 跳过注释和空行
        [[ "$key" =~ ^[[:space:]]*# ]] && continue
        [[ -z "$key" ]] && continue

        # 去除前后空格和引号
        key=$(echo "$key" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        value=$(echo "$value" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//;s/^["'\'']//;s/["'\'']$//')

        case "$key" in
            BACKUP_DIR)
                BACKUP_DIR="${BACKUP_DIR:-$value}"
                ;;
            MYSQL_HOST)
                MYSQL_HOST="${MYSQL_HOST:-$value}"
                ;;
            MYSQL_PORT)
                MYSQL_PORT="${MYSQL_PORT:-$value}"
                ;;
            MYSQL_USER)
                MYSQL_USER="${MYSQL_USER:-$value}"
                ;;
            MYSQL_PASSWORD)
                MYSQL_PASSWORD="${MYSQL_PASSWORD:-$value}"
                ;;
            DB_NAME)
                DB_NAME="${DB_NAME:-$value}"
                ;;
            REMAIN_DAY)
                REMAIN_DAY="${REMAIN_DAY:-$value}"
                ;;
        esac
    done < "$config_file"
}


load_config(){
      #优先命令行
      if [[ -n "$CONFIG_FILE" ]]; then
          load_from_config "$CONFIG_FILE"
      fi

      #否则从配置文件
      if [[ -z "$CONFIG_FILE" && -f "$DEFAULT_CONFIG_FILE" ]]; then
          load_from_config "$DEFAULT_CONFIG_FILE"
      fi

      #不行就环境变量
      load_from_env

      #最终设置
      BACKUP_DIR="${BACKUP_DIR:-$DEFAULT_BACKUP_DIR}"
      MYSQL_HOST="${MYSQL_HOST:-localhost}"
      MYSQL_PORT="${MYSQL_PORT:-3306}"
      REMAIN_DAY="${REMAIN_DAY:-$DEFAULT_REMAIN_DAY}"
}

#检查依赖
check() {
    local deps=("mysqldump" "gzip" "find")
    for dep in "${deps[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            echo "fail to find $dep" >&2
            exit 1
        fi
    done
}

#创建备份目录
create_backup_dir() {
    if ! mkdir -p "$BACKUP_DIR"; then
        echo "fail to create $BACKUP_DIR" >&2
        exit 1
    fi
}

backup(){
    loacl backup_file="$BACKUP_DIR/${DB_NAME}_$(date +%Y%m%d_%H%M%S).sql.gz"
    local tmep=$(mktemp)
    local flag=false

    #只有拥有者可以读写
    chmod 600 "$temp"

    if mysqldump --defaults-file="$temp" --single-transaction --quick \
      --databases "$DB_NAME" | gzip > "$backup_file";then
        if [[ -s "$backup_file" ]];then
          flag=true
        else
          rm -f "$backup_file"
        fi
    else
      rm -f "$backup_file"
    fi

    rm -f "$temp"

    if [[ "$flag" != "true" ]] then
        exit 1
    fi

}

clean() {
    local del_file
    del_file=$(find "$BACKUP_DIR" -name "${DB_NAME}_*.sql.gz" -type f -mtime "+$REMAIN_DAY" -delete -print)

    if [[ -n "$deleted_files" ]]; then
            echo "$deleted_files"
        else
            echo "no such files"
        fi
}



main() {
    parse "$@"
    load_config
    check
    create_backup_dir
    backup
    clean

}

main "$@"










