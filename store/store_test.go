package store

import (
	"context"
	"main/config"
	"main/testutil"
	"strconv"
	"strings"
	"testing"

	"github.com/go-redis/redis/v8"
)

func TestNewKVS(t *testing.T) {
	type args struct {
		ctx context.Context
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *KVS
		wantErr bool
	}{
		{name: "正常系 - 開発環境",args: args{ctx: context.Background(),cfg: &config.Config{Env: "dev" ,RedisHost: "127.0.0.1",RedisPort: 6379,RedisPassword: "",RedisTLS: ""}},want: &KVS{},wantErr:  false},
		{name: "正常系 - 本番環境",args: args{ctx: context.Background(),cfg: &config.Config{Env: "prod" ,RedisHost: "apn1-polished-sunbird-35172.upstash.io",RedisPort: 35172,RedisPassword: "d03c29220fe6412289943eac3d9aae1c",RedisTLS: "apn1-polished-sunbird-35172.upstash.io"}},want: &KVS{},wantErr:  false},
		{name: "異常系 - 開発環境",args: args{ctx: context.Background(),cfg: &config.Config{Env: "dev",RedisHost: "127.0.0.1",RedisPort: 36379,RedisPassword: "",RedisTLS: ""}},want: &KVS{},wantErr:  true},
		{name: "異常系 - 本番環境",args: args{ctx: context.Background(),cfg: &config.Config{Env: "prod",RedisHost: "apn1-polished-sunbird-35172.upstash.io",RedisPort: 35173,RedisPassword: "d03c29220fe6412289943eac3d9aae1c",RedisTLS: "apn1-polished-sunbird-35172.upstash.io"}},want: &KVS{},wantErr:  true},
	}

	for _, tt := range tests {	
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("GAME_ENV",tt.args.cfg.Env)
			t.Setenv("GAME_REDIS_HOST",tt.args.cfg.RedisHost)
			t.Setenv("GAME_REDIS_PORT",strconv.Itoa(tt.args.cfg.RedisPort))
			t.Setenv("GAME_REDIS_PASSWORD",tt.args.cfg.RedisPassword)
			t.Setenv("GAME_REDIS_TLS_SERVER_NAME",tt.args.cfg.RedisTLS)
	
			tt.want.Cli = testutil.OpenRedisForTest(t)
			got, err := NewKVS(tt.args.ctx, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKVS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(tt.name,"正常系") {return}
			if strings.Compare(got.Cli.Options().Addr,tt.want.Cli.Options().Addr) != 0 {
				t.Errorf("NewKVS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKVS_Save(t *testing.T) {
	type fields struct {
		Cli *redis.Client
	}
	type args struct {
		ctx   context.Context
		key   string
		value int
	}
	tests := []struct {
		name    string
		cfg     *config.Config
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "正常系 - 開発環境",cfg: &config.Config{Env: "dev" ,RedisHost: "127.0.0.1",RedisPort: 6379,RedisPassword: "",RedisTLS: ""},fields: fields{Cli: &redis.Client{}},args: args{ctx: context.Background(),key: "hoge",value:100},wantErr: false},
		{name: "正常系 - 本番環境",cfg: &config.Config{Env: "prod" ,RedisHost: "apn1-polished-sunbird-35172.upstash.io",RedisPort: 35172,RedisPassword: "d03c29220fe6412289943eac3d9aae1c",RedisTLS: "apn1-polished-sunbird-35172.upstash.io"},fields: fields{Cli: &redis.Client{}},args: args{ctx: context.Background(),key: "hoge",value:100},wantErr: false},
		{name: "異常系 - 開発環境",cfg: &config.Config{Env: "dev" ,RedisHost: "127.0.0.1",RedisPort: 6370,RedisPassword: "",RedisTLS: ""},fields: fields{Cli: &redis.Client{}},args: args{ctx: context.Background(),key: "hoge",value:100},wantErr: true},
		{name: "異常系 - 本番環境",cfg: &config.Config{Env: "prod" ,RedisHost: "apn1-polished-sunbird-35172.upstash.io",RedisPort: 6379,RedisPassword: "",RedisTLS: ""},fields: fields{Cli: &redis.Client{}},args: args{ctx: context.Background(),key: "hoge",value:100},wantErr: true},
	}
	for _, tt := range tests {
		t.Setenv("GAME_ENV",tt.cfg.Env)
		t.Setenv("GAME_REDIS_HOST",tt.cfg.RedisHost)
		t.Setenv("GAME_REDIS_PORT",strconv.Itoa(tt.cfg.RedisPort))
		t.Setenv("GAME_REDIS_PASSWORD",tt.cfg.RedisPassword)
		t.Setenv("GAME_REDIS_TLS_SERVER_NAME",tt.cfg.RedisTLS)

		t.Run(tt.name, func(t *testing.T) {
			k := &KVS{
				Cli: testutil.OpenRedisForTest(t),
			}
			if !strings.Contains(tt.name,"正常系") {return}
			if err := k.Save(tt.args.ctx, tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("KVS.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKVS_Load(t *testing.T) {
	type fields struct {
		Cli *redis.Client
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		cfg     *config.Config
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{name: "正常系 - 開発環境",cfg: &config.Config{Env: "dev" ,RedisHost: "127.0.0.1",RedisPort: 6379,RedisPassword: "",RedisTLS: ""},fields: fields{Cli: &redis.Client{}},args: args{ctx: context.Background(),key: "hoge"},want:100,wantErr: false},
		{name: "正常系 - 本番環境",cfg: &config.Config{Env: "prod" ,RedisHost: "apn1-polished-sunbird-35172.upstash.io",RedisPort: 35172,RedisPassword: "d03c29220fe6412289943eac3d9aae1c",RedisTLS: "apn1-polished-sunbird-35172.upstash.io"},fields: fields{Cli: &redis.Client{}},args: args{ctx: context.Background(),key: "hoge"},want:100,wantErr: false},
		{name: "異常系 - 開発環境",cfg: &config.Config{Env: "dev" ,RedisHost: "127.0.0.1",RedisPort: 6370,RedisPassword: "",RedisTLS: ""},fields: fields{Cli: &redis.Client{}},args: args{ctx: context.Background(),key: "hoge"},want:100,wantErr: true},
		{name: "異常系 - 本番環境",cfg: &config.Config{Env: "prod" ,RedisHost: "apn1-polished-sunbird-35172.upstash.io",RedisPort: 6379,RedisPassword: "",RedisTLS: ""},fields: fields{Cli: &redis.Client{}},args: args{ctx: context.Background(),key: "hoge"},want:100,wantErr: true},
	}
	for _, tt := range tests {
		t.Setenv("GAME_ENV",tt.cfg.Env)
		t.Setenv("GAME_REDIS_HOST",tt.cfg.RedisHost)
		t.Setenv("GAME_REDIS_PORT",strconv.Itoa(tt.cfg.RedisPort))
		t.Setenv("GAME_REDIS_PASSWORD",tt.cfg.RedisPassword)
		t.Setenv("GAME_REDIS_TLS_SERVER_NAME",tt.cfg.RedisTLS)

		t.Run(tt.name, func(t *testing.T) {
			k := &KVS{
				Cli: testutil.OpenRedisForTest(t),
			}
			got, err := k.Load(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("KVS.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(tt.name,"正常系") {return}
			if got != tt.want {
				t.Errorf("KVS.Load() = %v, want %v", got, tt.want)
			}
		})
	}
}
