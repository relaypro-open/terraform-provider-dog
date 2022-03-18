#resource "dog_profile" "new_profile" {
#  name = "new_name"
#  description = "test_profile"
#  rules = {
#    inbound = [
#      {
#        action = "ACCEPT"
#        active = true
#        comment = "accept all"
#        environments = []
#        group = "any"
#        group_type = "ANY"
#        interface = ""
#        log = false
#        log_prefix = ""
#        order = 1 
#        service = "any"
#        states = []
#        type = "BASIC"
#      }
#    ]
#    outbound = [
#      {
#        action = "ACCEPT"
#        active = true
#        comment = ""
#        environments = []
#        group = "any"
#        group_type = "ANY"
#        interface = ""
#        log = false
#        log_prefix = ""
#        order = 1
#        service = "any"
#        states = []
#        type = "BASIC"
#      }
#    ]
#  }
#  version = "1.0"
#}

//resource "dog_profile" "new_profile" {
//  name = "new_name"
//  description = "test_profile"
//  rules = {
//    inbound = [
//      {
//        action = "ACCEPT"
//        active = true
//        comment = "accept all"
//        environments = []
//        group = "any"
//        group_type = "ANY"
//        interface = ""
//        log = false
//        log_prefix = ""
//        order = 1 
//        service = "any"
//        states = []
//        type = "BASIC"
//      },
//      {
//        action = "DROP"
//        active = true
//        comment = "drop all"
//        environments = []
//        group = "any"
//        group_type = "ANY"
//        interface = ""
//        log = false
//        log_prefix = ""
//        order = 2
//        service = "any"
//        states = []
//        type = "BASIC"
//      }
//    ]
//    outbound = [
//      {
//        action = "ACCEPT"
//        active = true
//        comment = ""
//        environments = []
//        group = "any"
//        group_type = "ANY"
//        interface = ""
//        log = false
//        log_prefix = ""
//        order = 1
//        service = "any"
//        states = []
//        type = "BASIC"
//      }
//    ]
//  }
//  version = "1.0"
//}
