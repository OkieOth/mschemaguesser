db = db.getSiblingDB('dummy');
db.c1.insert({ key: 'value1', number: 12, bool: true });
db.c1.insert({ key: 'value2', number: 13, bool: true });
db.c1.insert({ key: 'value3', number: 14 });
db.c1.insert({ key: 'value4', number: 15, bool: true });

db.c2.insert({ complex: { name: "homer", array: [1, 2, 3, 4] }, bool: false });
db.c2.insert({ complex: { name: "marge", array: [1, 2, 3] }, bool: true });
db.c2.insert({ complex: { name: "maggy", array: [1, 2, 3], hobbies: { saxophon: true, skating: false } }, bool: true });


db.c3.insert({
      Metadata: {
        Created: new Date("2024-01-09T08:17:15.831Z"),
        Modified: new Date("2024-02-08T13:09:47.559Z"),
        Revision: 3,
        SchemaVersion: "1",
        State: 0
      },
      Resource: {
        CompositeThingKeys: ["yyy", "xxx"],
        CompositeThingsMap: {
          yyy: {
            SimpleThingKeys: ["yyy_category", "yyy_situation", "yyy_localStrategy"],
            SimpleThingsMap: {
              yyy_category: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "2" }
              },
              yyy_localStrategy: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "3" }
              },
              yyy_situation: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "1" }
              }
            }
          },
          xxx: {
            SimpleThingKeys: ["xxx_schnulli", "xxx_nodeState", "xxx_individualTrafficDependentModification"],
            SimpleThingsMap: {
              xxx_individualTrafficDependentModification: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "3" }
              },
              xxx_nodeState: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "1" }
              },
              xxx_schnulli: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "1" }
              }
            }
          }
        },
        Description: "fasfasöf asöd fasöldfösadfösdf",
        Guid: new BinData(4, "5c46d5a7f28e45578c3e6221d08552b6"),
        IsTemplate: true,
        Name: "a name",
        OrgId: new BinData(4, "48d387c4f3c14b308dc1a91dec1fe7e5")
      }
});
db.c3.insert({
      Metadata: {
        Created: new Date("2024-01-09T16:00:02.55Z"),
        Modified: new Date("2024-01-09T16:00:02.569Z"),
        Revision: 0,
        SchemaVersion: "1",
        State: 0
      },
      Resource: {
        CompositeThingKeys: ["yyy"],
        CompositeThingsMap: {
          yyy: {
            SimpleThingKeys: ["yyy_localStrategy"],
            SimpleThingsMap: {
              yyy_localStrategy: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "2" }
              }
            }
          }
        },
        Description: "",
        Guid: new BinData(4, "ace326c30361464b8c682ec694acb516"),
        IsTemplate: true,
        Name: "Local2",
        OrgId: new BinData(4, "233b364b677e4f6dafcea43b53bdb126")
      }
});
db.c3.insert({
      Metadata: {
        Created: new Date("2024-01-09T16:00:15.551Z"),
        Modified: new Date("2024-01-09T16:00:15.569Z"),
        Revision: 0,
        SchemaVersion: "1",
        State: 0
      },
      Resource: {
        CompositeThingKeys: ["yyy"],
        CompositeThingsMap: {
          yyy: {
            SimpleThingKeys: ["yyy_local"],
            SimpleThingsMap: {
              yyy_local: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "3" }
              }
            }
          }
        },
        Description: "",
        Guid: new BinData(4, "07bfbe1aa8bf42c2bea246df94b5ad87"),
        IsTemplate: true,
        Name: "Local3",
        OrgId: new BinData(4, "233b364b677e4f6dafcea43b53bdb126")
      }
});
db.c3.insert({
      Metadata: {
        Created: new Date("2024-01-10T08:58:22.139Z"),
        Modified: new Date("2024-01-10T08:58:22.158Z"),
        Revision: 0,
        SchemaVersion: "1",
        State: 0
      },
      Resource: {
        CompositeThingKeys: ["yyy"],
        CompositeThingsMap: {
          yyy: {
            SimpleThingKeys: ["yyy_localStrategy"],
            SimpleThingsMap: {
              yyy_localStrategy: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "4" }
              }
            }
          }
        },
        Description: "",
        Guid: new BinData(4, "853644d0a9f044ad8cdfb6dc9a7845ed"),
        IsTemplate: true,
        Name: "Local4",
        OrgId: new BinData(4, "233b364b677e4f6dafcea43b53bdb126")
      }
});
db.c3.insert({
      Metadata: {
        Created: new Date("2024-01-10T09:47:00.599Z"),
        Modified: new Date("2024-01-23T15:04:25.812Z"),
        Revision: 2,
        SchemaVersion: "1",
        State: 0
      },
      Resource: {
        CompositeThingKeys: ["yyy"],
        CompositeThingsMap: {
          yyy: {
            SimpleThingKeys: ["yyy_category"],
            SimpleThingsMap: {
              yyy_category: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "1" }
              }
            }
          }
        },
        Description: "Test command",
        Guid: new BinData(4, "0b7835117bd840a6b9e6cf34a4607885"),
        IsTemplate: true,
        Name: "Test name",
        OrgId: new BinData(4, "a5dddb893b9d4ae58340888f2d9e2dc8")
      }
});
db.c3.insert({
      Metadata: {
        Created: new ISODate("2024-01-11T07:38:00.397Z"),
        Modified: new ISODate("2024-01-11T07:38:01.071Z"),
        Revision: 0,
        SchemaVersion: "1",
        State: 0
      },
      Resource: {
        CompositeThingKeys: ["xxx"],
        CompositeThingsMap: {
          xxx: {
            SimpleThingKeys: ["xxx_nodeState"],
            SimpleThingsMap: {
              xxx_nodeState: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "enabled" }
              }
            }
          }
        },
        Description: "NodeState Enabled",
        Guid: new UUID("ce523ea0-f45b-438a-b135-199530eef6ee"),
        IsTemplate: true,
        Name: "NS_Enabled",
        OrgId: new UUID("48d387c4-f3c1-4b30-8dc1-a91dec1fe7e5")
      }
    });
db.c3.insert({
      Metadata: {
        Created: new ISODate("2024-01-11T13:08:32.982Z"),
        Modified: new ISODate("2024-02-02T17:07:58.461Z"),
        Revision: 3,
        SchemaVersion: "1",
        State: 0
      },
      Resource: {
        CompositeThingKeys: ["xxx"],
        CompositeThingsMap: {
          xxx: {
            SimpleThingKeys: ["xxx_schnulli", "xxx_nodeState", "xxx_intervention"],
            SimpleThingsMap: {
              xxx_intervention: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "0" }
              },
              xxx_nodeState: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "ENABLED" }
              },
              xxx_schnulli: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "1" }
              }
            }
          }
        },
        Description: "ThingXXXX",
        Guid: new UUID("718981e4-e0f1-4433-99bf-d9d924dccf7a"),
        IsTemplate: true,
        Name: "P1, Enabled",
        OrgId: new UUID("a5dddb89-3b9d-4ae5-8340-888f2d9e2dc8")
      }
    });
db.c3.insert({
      Metadata: {
        Created: new ISODate("2024-01-15T08:55:39.733Z"),
        Modified: new ISODate("2024-01-15T12:23:00.649Z"),
        Revision: 1,
        SchemaVersion: "1",
        State: 0
      },
      Resource: {
        CompositeThingKeys: ["xxx", "yyy"],
        CompositeThingsMap: {
          yyy: {
            SimpleThingKeys: ["yyy_localStrategy"],
            SimpleThingsMap: {
              yyy_localStrategy: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "1" }
              }
            }
          },
          xxx: {
            SimpleThingKeys: ["xxx_schnulli"],
            SimpleThingsMap: {
              xxx_schnulli: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "1" }
              }
            }
          }
        },
        Description: "",
        Guid: new UUID("a20a1055-8beb-447a-9746-9a8afde87fb4"),
        IsTemplate: true,
        Name: "S01",
        OrgId: new UUID("a5dddb89-3b9d-4ae5-8340-888f2d9e2dc8")
      }
    });
db.c3.insert({
      Metadata: {
        Created: new ISODate("2024-01-15T08:55:52.092Z"),
        Modified: new ISODate("2024-01-15T12:23:13.663Z"),
        Revision: 1,
        SchemaVersion: "1",
        State: 0
      },
      Resource: {
        CompositeThingKeys: ["xxx", "yyy"],
        CompositeThingsMap: {
          yyy: {
            SimpleThingKeys: ["yyy_localStrategy"],
            SimpleThingsMap: {
              yyy_localStrategy: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "2" }
              }
            }
          },
          xxx: {
            SimpleThingKeys: ["xxx_schnulli"],
            SimpleThingsMap: {
              xxx_schnulli: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "2" }
              }
            }
          }
        },
        Description: "",
        Guid: new UUID("69101d92-4b16-49b8-bc7a-410dd33c3cee"),
        IsTemplate: true,
        Name: "S02",
        OrgId: new UUID("a5dddb89-3b9d-4ae5-8340-888f2d9e2dc8")
      }
    });
db.c3.insert({
      Metadata: {
        Created: new ISODate("2024-01-15T12:17:35.166Z"),
        Modified: new ISODate("2024-02-05T10:34:12.848Z"),
        Revision: 3,
        SchemaVersion: "1",
        State: 0
      },
      Resource: {
        CompositeThingKeys: ["yyy", "xxx"],
        CompositeThingsMap: {
          yyy: {
            SimpleThingKeys: ["yyy_localStrategy", "yyy_busPriority"],
            SimpleThingsMap: {
              yyy_busPriority: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "12" }
              },
              yyy_localStrategy: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "12" }
              }
            }
          },
          xxx: {
            SimpleThingKeys: ["xxx_schnulli", "xxx_publicTransport"],
            SimpleThingsMap: {
              xxx_publicTransport: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "21" }
              },
              xxx_schnulli: {
                ParameterKeys: ["value"],
                ParametersMap: { value: "3" }
              }
            }
          }
        },
        Description: "S03",
        Guid: new UUID("be783abf-413c-4c2c-a80c-8752f0436f90"),
        IsTemplate: true,
        Name: "S03",
        OrgId: new UUID("a5dddb89-3b9d-4ae5-8340-888f2d9e2dc8")
      }
    });
db.c3.insert({
    Metadata: {
    Created: new ISODate("2024-01-22T10:25:54.166Z"),
    Modified: new ISODate("2024-01-23T15:43:08.959Z"),
    Revision: 4,
    SchemaVersion: "1",
    State: 0
    },
    Resource: {
    CompositeThingKeys: ["xxx"],
    CompositeThingsMap: {
        xxx: {
        SimpleThingKeys: ["xxx_schnulli", "xxx_nodeState", "xxx_subNodeState-0"],
        SimpleThingsMap: {
            xxx_nodeState: {
            ParameterKeys: ["value"],
            ParametersMap: { value: "ENABLED" }
            },
            xxx_schnulli: {
            ParameterKeys: ["value"],
            ParametersMap: { value: "1" }
            },
            xxx_subNodeState: {
            ParameterKeys: ["value"],
            ParametersMap: { value: "ENABLED" }
            }
        }
        }
    },
    Description: "",
    Guid: new UUID("ac7e7a94-165b-4d80-a5c6-365e3b294c0e"),
    IsTemplate: true,
    Name: "test P1",
    OrgId: new UUID("a5dddb89-3b9d-4ae5-8340-888f2d9e2dc8")
    }
});
db.c3.insert({
    Metadata: {
    Created: new ISODate("2024-01-22T10:30:11.81Z"),
    Modified: new ISODate("2024-02-06T13:04:17.801Z"),
    Revision: 4,
    SchemaVersion: "1",
    State: 0
    },
    Resource: {
    CompositeThingKeys: ["xxx"],
    CompositeThingsMap: {
        xxx: {
        SimpleThingKeys: ["xxx_schnulli", "xxx_nodeState"],
        SimpleThingsMap: {
            xxx_nodeState: {
            ParameterKeys: ["value"],
            ParametersMap: { value: "ENABLED" }
            },
            xxx_schnulli: {
            ParameterKeys: ["value"],
            ParametersMap: { value: "2" }
            }
        }
        }
    },
    Description: "",
    Guid: new UUID("150aba60-869a-4d05-9cd1-fcf396d055b6"),
    IsTemplate: true,
    Name: "P2",
    OrgId: new UUID("a5dddb89-3b9d-4ae5-8340-888f2d9e2dc8")
    }
});
db.c3.insert({
    Metadata: {
    Created: new ISODate("2024-01-22T11:13:47.859Z"),
    Modified: new ISODate("2024-01-22T11:13:47.939Z"),
    Revision: 0,
    SchemaVersion: "1",
    State: 0
    },
    Resource: {
    CompositeThingKeys: ["xxx"],
    CompositeThingsMap: {
        xxx: {
        SimpleThingKeys: ["xxx_schnulli", "xxx_nodeState"],
        SimpleThingsMap: {
            xxx_nodeState: {
            ParameterKeys: ["value"],
            ParametersMap: { value: "ENABLED" }
            },
            xxx_schnulli: {
            ParameterKeys: ["value"],
            ParametersMap: { value: "8" }
            }
        }
        }
    },
    Description: "",
    Guid: new UUID("2ed9fac2-dd63-44d3-9807-1aad6c69fc64"),
    IsTemplate: true,
    Name: "p8",
    OrgId: new UUID("a5dddb89-3b9d-4ae5-8340-888f2d9e2dc8")
    }
});
db.c3.insert({
    Metadata: {
    Created: new ISODate("2024-01-22T13:01:28.836Z"),
    Modified: new ISODate("2024-01-22T13:01:28.899Z"),
    Revision: 0,
    SchemaVersion: "1",
    State: 0
    },
    Resource: {
    CompositeThingKeys: ["yyy"],
    CompositeThingsMap: {
        yyy: {
        SimpleThingKeys: ["yyy_localStrategy"],
        SimpleThingsMap: {
            yyy_localStrategy: {
            ParameterKeys: ["value"],
            ParametersMap: { value: "12" }
            }
        }
        }
    },
    Description: "Globa Thing C",
    Guid: new UUID("ae4f5010-f315-48fb-8a47-18bbc42f0984"),
    IsTemplate: true,
    Name: "Global afasdfa",
    OrgId: new UUID("a5dddb89-3b9d-4ae5-8340-888f2d9e2dc8")
    }
});
db.c3.insert({
    Metadata: {
    Created: new ISODate("2024-01-22T17:50:12.508Z"),
    Modified: new ISODate("2024-02-16T08:54:18.055Z"),
    Revision: 3,
    SchemaVersion: "1",
    State: 0
    },
    Resource: {
    CompositeThingKeys: ["xxx"],
    CompositeThingsMap: {
        xxx: {
        SimpleThingKeys: ["xxx_schnulli", "xxx_projectModification-0"],
        SimpleThingsMap: {
            xxx_projectModification0: {
            ParameterKeys: ["value"],
            ParametersMap: { value: "2" }
            },
            xxx_schnulli: {
            ParameterKeys: ["value"],
            ParametersMap: { value: "1" }
            }
        }
        }
    },
    Description: "test",
    Guid: new UUID("275a96a9-27f7-4c2d-83b9-34162a112d31"),
    IsTemplate: true,
    Name: "test-command",
    OrgId: new UUID("48d387c4-f3c1-4b30-8dc1-a91dec1fe7e5")
    }
});