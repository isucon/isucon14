<?php

/**
 * ISURIDE API Specification
 * PHP version 7.4
 *
 * @package IsuRide
 * @author  OpenAPI Generator team
 * @link    https://github.com/openapitools/openapi-generator
 */

/**
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 * The version of the OpenAPI document: 1.0
 * Generated by: https://github.com/openapitools/openapi-generator.git
 */

declare(strict_types=1);

/**
 * NOTE: This class is auto generated by the openapi generator program.
 * https://github.com/openapitools/openapi-generator
 * Do not edit the class manually.
 */
namespace IsuRide\App;

use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use Slim\Exception\HttpNotImplementedException;

/**
 * RegisterRoutes Class Doc Comment
 *
 * @package IsuRide
 * @author  OpenAPI Generator team
 * @link    https://github.com/openapitools/openapi-generator
 */
class RegisterRoutes
{
    /** @var array[] list of all api operations */
    private $operations = [
        [
            'httpMethod' => 'GET',
            'basePathWithoutHost' => '/api',
            'path' => '/app/nearby-chairs',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractAppApi',
            'userClassname' => 'AppApi',
            'operationId' => 'appGetNearbyChairs',
            'responses' => [
                '200' => [
                    'jsonSchema' => '{
  "description" : "OK",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/app_get_nearby_chairs_200_response"
      }
    }
  }
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'GET',
            'basePathWithoutHost' => '/api',
            'path' => '/app/notification',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractAppApi',
            'userClassname' => 'AppApi',
            'operationId' => 'appGetNotification',
            'responses' => [
                '200' => [
                    'jsonSchema' => '{
  "description" : "OK",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/app_get_notification_200_response"
      }
    }
  }
}',
                ],
                '204' => [
                    'jsonSchema' => '{
  "description" : "対象となるライドが存在しない場合"
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'GET',
            'basePathWithoutHost' => '/api',
            'path' => '/app/rides',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractAppApi',
            'userClassname' => 'AppApi',
            'operationId' => 'appGetRides',
            'responses' => [
                '200' => [
                    'jsonSchema' => '{
  "description" : "OK",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/app_get_rides_200_response"
      }
    }
  }
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'POST',
            'basePathWithoutHost' => '/api',
            'path' => '/app/payment-methods',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractAppApi',
            'userClassname' => 'AppApi',
            'operationId' => 'appPostPaymentMethods',
            'responses' => [
                '204' => [
                    'jsonSchema' => '{
  "description" : "決済トークンの登録に成功した"
}',
                ],
                '400' => [
                    'jsonSchema' => '{
  "description" : "Bad Request",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/Error"
      }
    }
  }
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'POST',
            'basePathWithoutHost' => '/api',
            'path' => '/app/rides',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractAppApi',
            'userClassname' => 'AppApi',
            'operationId' => 'appPostRides',
            'responses' => [
                '202' => [
                    'jsonSchema' => '{
  "description" : "配車要求を受け付けた",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/app_post_rides_202_response"
      }
    }
  }
}',
                ],
                '400' => [
                    'jsonSchema' => '{
  "description" : "Bad Request",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/Error"
      }
    }
  }
}',
                ],
                '409' => [
                    'jsonSchema' => '{
  "description" : "Conflict",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/Error"
      }
    }
  }
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'POST',
            'basePathWithoutHost' => '/api',
            'path' => '/app/rides/estimated-fare',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractAppApi',
            'userClassname' => 'AppApi',
            'operationId' => 'appPostRidesEstimatedFare',
            'responses' => [
                '200' => [
                    'jsonSchema' => '{
  "description" : "OK",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/app_post_rides_estimated_fare_200_response"
      }
    }
  }
}',
                ],
                '400' => [
                    'jsonSchema' => '{
  "description" : "Bad Request",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/Error"
      }
    }
  }
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'POST',
            'basePathWithoutHost' => '/api',
            'path' => '/app/users',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractAppApi',
            'userClassname' => 'AppApi',
            'operationId' => 'appPostUsers',
            'responses' => [
                '201' => [
                    'jsonSchema' => '{
  "description" : "ユーザー登録が完了した",
  "headers" : {
    "Set-Cookie" : {
      "description" : "サーバーから返却される Cookie",
      "style" : "simple",
      "explode" : false,
      "schema" : {
        "type" : "string",
        "example" : "app_session=<access_token>; Path=/;"
      }
    }
  },
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/app_post_users_201_response"
      }
    }
  }
}',
                ],
                '400' => [
                    'jsonSchema' => '{
  "description" : "Bad Request",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/Error"
      }
    }
  }
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'GET',
            'basePathWithoutHost' => '/api',
            'path' => '/app/rides/{ride_id}',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractAppApi',
            'userClassname' => 'AppApi',
            'operationId' => 'appGetRide',
            'responses' => [
                '200' => [
                    'jsonSchema' => '{
  "description" : "OK",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/AppRide"
      }
    }
  }
}',
                ],
                '404' => [
                    'jsonSchema' => '{
  "description" : "Not Found",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/Error"
      }
    }
  }
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'POST',
            'basePathWithoutHost' => '/api',
            'path' => '/app/rides/{ride_id}/evaluation',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractAppApi',
            'userClassname' => 'AppApi',
            'operationId' => 'appPostRideEvaluation',
            'responses' => [
                '200' => [
                    'jsonSchema' => '{
  "description" : "ユーザーがライドを評価した",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/app_post_ride_evaluation_200_response"
      }
    }
  }
}',
                ],
                '400' => [
                    'jsonSchema' => '{
  "description" : "椅子が目的地に到着していない、ユーザーが乗車していない、すでに到着しているなど",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/Error"
      }
    }
  }
}',
                ],
                '404' => [
                    'jsonSchema' => '{
  "description" : "存在しないライド",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/Error"
      }
    }
  }
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'GET',
            'basePathWithoutHost' => '/api',
            'path' => '/chair/notification',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractChairApi',
            'userClassname' => 'ChairApi',
            'operationId' => 'chairGetNotification',
            'responses' => [
                '200' => [
                    'jsonSchema' => '{
  "description" : "自分に割り当てられた最新のライド",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/chair_get_notification_200_response"
      }
    }
  }
}',
                ],
                '204' => [
                    'jsonSchema' => '{
  "description" : "割り当てられたライドが存在しない場合"
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'POST',
            'basePathWithoutHost' => '/api',
            'path' => '/chair/activity',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractChairApi',
            'userClassname' => 'ChairApi',
            'operationId' => 'chairPostActivity',
            'responses' => [
                '204' => [
                    'jsonSchema' => '{
  "description" : "椅子の配車受付の開始・停止を受理した"
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'POST',
            'basePathWithoutHost' => '/api',
            'path' => '/chair/chairs',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractChairApi',
            'userClassname' => 'ChairApi',
            'operationId' => 'chairPostChairs',
            'responses' => [
                '201' => [
                    'jsonSchema' => '{
  "description" : "椅子登録が完了した",
  "headers" : {
    "Set-Cookie" : {
      "description" : "サーバーから返却される Cookie",
      "style" : "simple",
      "explode" : false,
      "schema" : {
        "type" : "string",
        "example" : "chair_session=<access_token>; Path=/;"
      }
    }
  },
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/chair_post_chairs_201_response"
      }
    }
  }
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'POST',
            'basePathWithoutHost' => '/api',
            'path' => '/chair/coordinate',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractChairApi',
            'userClassname' => 'ChairApi',
            'operationId' => 'chairPostCoordinate',
            'responses' => [
                '200' => [
                    'jsonSchema' => '{
  "description" : "椅子の座標を更新した",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/chair_post_coordinate_200_response"
      }
    }
  }
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'GET',
            'basePathWithoutHost' => '/api',
            'path' => '/chair/rides/{ride_id}',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractChairApi',
            'userClassname' => 'ChairApi',
            'operationId' => 'chairGetRide',
            'responses' => [
                '200' => [
                    'jsonSchema' => '{
  "description" : "OK",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/ChairRide"
      }
    }
  }
}',
                ],
                '404' => [
                    'jsonSchema' => '{
  "description" : "存在しないライド、対象の椅子にマッチングされていないライドを取得しようとした場合など",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/Error"
      }
    }
  }
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'POST',
            'basePathWithoutHost' => '/api',
            'path' => '/chair/rides/{ride_id}/status',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractChairApi',
            'userClassname' => 'ChairApi',
            'operationId' => 'chairPostRideStatus',
            'responses' => [
                '204' => [
                    'jsonSchema' => '{
  "description" : "No Content"
}',
                ],
                '404' => [
                    'jsonSchema' => '{
  "description" : "Not Found",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/Error"
      }
    }
  }
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'GET',
            'basePathWithoutHost' => '/api',
            'path' => '/owner/chairs',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractOwnerApi',
            'userClassname' => 'OwnerApi',
            'operationId' => 'ownerGetChairs',
            'responses' => [
                '200' => [
                    'jsonSchema' => '{
  "description" : "OK",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/owner_get_chairs_200_response"
      }
    }
  }
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'GET',
            'basePathWithoutHost' => '/api',
            'path' => '/owner/sales',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractOwnerApi',
            'userClassname' => 'OwnerApi',
            'operationId' => 'ownerGetSales',
            'responses' => [
                '200' => [
                    'jsonSchema' => '{
  "description" : "OK",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/owner_get_sales_200_response"
      }
    }
  }
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'POST',
            'basePathWithoutHost' => '/api',
            'path' => '/owner/owners',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractOwnerApi',
            'userClassname' => 'OwnerApi',
            'operationId' => 'ownerPostOwners',
            'responses' => [
                '201' => [
                    'jsonSchema' => '{
  "description" : "オーナー登録が完了した",
  "headers" : {
    "Set-Cookie" : {
      "description" : "サーバーから返却される Cookie",
      "style" : "simple",
      "explode" : false,
      "schema" : {
        "type" : "string",
        "example" : "owner_session=<access_token>; Path=/;"
      }
    }
  },
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/owner_post_owners_201_response"
      }
    }
  }
}',
                ],
                '400' => [
                    'jsonSchema' => '{
  "description" : "Bad Request",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/Error"
      }
    }
  }
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'GET',
            'basePathWithoutHost' => '/api',
            'path' => '/owner/chairs/{chair_id}',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractOwnerApi',
            'userClassname' => 'OwnerApi',
            'operationId' => 'ownerGetChair',
            'responses' => [
                '200' => [
                    'jsonSchema' => '{
  "description" : "OK",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/owner_get_chairs_200_response_chairs_inner"
      }
    }
  }
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
        [
            'httpMethod' => 'POST',
            'basePathWithoutHost' => '/api',
            'path' => '/initialize',
            'apiPackage' => 'IsuRide\Api',
            'classname' => 'AbstractSystemApi',
            'userClassname' => 'SystemApi',
            'operationId' => 'postInitialize',
            'responses' => [
                '200' => [
                    'jsonSchema' => '{
  "description" : "サービスの初期化が完了した",
  "content" : {
    "application/json" : {
      "schema" : {
        "$ref" : "#/components/schemas/post_initialize_200_response"
      }
    }
  }
}',
                ],
            ],
            'authMethods' => [
            ],
        ],
    ];

    /**
     * Add routes to Slim app.
     *
     * @param \Slim\App $app Pre-configured Slim application instance
     *
     * @throws HttpNotImplementedException When implementation class doesn't exists
     */
    public function __invoke(\Slim\App $app): void
    {
        $app->options('/{routes:.*}', function (ServerRequestInterface $request, ResponseInterface $response) {
            // CORS Pre-Flight OPTIONS Request Handler
            return $response;
        });

        // create mock middleware factory
        /** @var \Psr\Container\ContainerInterface */
        $container = $app->getContainer();
        /** @var \OpenAPIServer\Mock\OpenApiDataMockerRouteMiddlewareFactory|null */
        $mockMiddlewareFactory = null;
        if ($container->has(\OpenAPIServer\Mock\OpenApiDataMockerRouteMiddlewareFactory::class)) {
            // I know, anti-pattern. Don't retrieve dependency directly from container
            $mockMiddlewareFactory = $container->get(\OpenAPIServer\Mock\OpenApiDataMockerRouteMiddlewareFactory::class);
        }

        foreach ($this->operations as $operation) {
            $callback = function (ServerRequestInterface $request) use ($operation) {
                $message = "How about extending {$operation['classname']} by {$operation['apiPackage']}\\{$operation['userClassname']} class implementing {$operation['operationId']} as a {$operation['httpMethod']} method?";
                throw new HttpNotImplementedException($request, $message);
            };
            $middlewares = [];

            if (class_exists("\\{$operation['apiPackage']}\\{$operation['userClassname']}")) {
                // Notice how we register the controller using the class name?
                // PHP-DI will instantiate the class for us only when it's actually necessary
                $callback = ["\\{$operation['apiPackage']}\\{$operation['userClassname']}", $operation['operationId']];
            }

            if ($mockMiddlewareFactory) {
                $mockSchemaResponses = array_map(function ($item) {
                    return json_decode($item['jsonSchema'], true);
                }, $operation['responses']);
                $middlewares[] = $mockMiddlewareFactory->create($mockSchemaResponses);
            }

            $route = $app->map(
                [$operation['httpMethod']],
                "{$operation['basePathWithoutHost']}{$operation['path']}",
                $callback
            )->setName($operation['operationId']);


            foreach ($middlewares as $middleware) {
                $route->add($middleware);
            }
        }
    }
}