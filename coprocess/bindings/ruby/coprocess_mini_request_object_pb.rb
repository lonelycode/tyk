# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: coprocess_mini_request_object.proto

require 'google/protobuf'

require 'coprocess_return_overrides_pb'
Google::Protobuf::DescriptorPool.generated_pool.build do
  add_message "coprocess.MiniRequestObject" do
    map :headers, :string, :string, 1
    map :set_headers, :string, :string, 2
    repeated :delete_headers, :string, 3
    optional :body, :string, 4
    optional :url, :string, 5
    map :params, :string, :string, 6
    map :add_params, :string, :string, 7
    map :extended_params, :string, :string, 8
    repeated :delete_params, :string, 9
    optional :return_overrides, :message, 10, "coprocess.ReturnOverrides"
    optional :method, :string, 11
    optional :request_uri, :string, 12
    optional :scheme, :string, 13
    optional :raw_body, :bytes, 14
  end
end

module Coprocess
  MiniRequestObject = Google::Protobuf::DescriptorPool.generated_pool.lookup("coprocess.MiniRequestObject").msgclass
end
