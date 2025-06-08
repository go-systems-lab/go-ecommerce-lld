import grpc
from concurrent import futures
from generated.pb import recommender_pb2_grpc, recommender_pb2
from app.services.recommender import recommender
from app.db.models import Product
from app.db.session import get_db_session


def _fetch_grpc_products(product_ids):
    """Helper method to fetch products from the DB and return gRPC ProductReplica objects."""
    with get_db_session() as session:
        products = (
            session.query(Product)
            .filter(Product.id.in_(product_ids))
            .all()
        )

    grpc_products = [
        recommender_pb2.ProductReplica(
            id=product.id,
            name=product.name,
            description=product.description,
            price=product.price,
        )
        for product in products
    ]
    return grpc_products


def _handle_exception(context, error_message):
    """Helper method to set gRPC status code and details."""
    context.set_code(grpc.StatusCode.INTERNAL)
    context.set_details(error_message)



class RecommenderServiceServicer(recommender_pb2_grpc.RecommenderServiceServicer):
    def GetRecommendationsForUserId(self, request, context):
        user_id = request.user_id
        skip = request.skip or 0
        take = request.take or 5

        try:
            recommended_product_ids = recommender.recommend_on_user_id(
                user_id=user_id,
                skip=skip,
                take=take
            )

            # Fetch product details as gRPC objects
            grpc_products = _fetch_grpc_products(recommended_product_ids)
            return recommender_pb2.RecommendationResponse(
                recommended_products=grpc_products
            )
        except Exception as e:
            _handle_exception(context, f"Error getting recommendations: {str(e)}")
            return recommender_pb2.RecommendationResponse()
        
    def GetRecommendationsOnViews(self, request, context):
        viewed_product_ids = request.ids
        skip = request.skip or 0
        take = request.take or 5

        try:
            recommended_product_ids = recommender.recommend_on_viewed_ids(
                viewed_ids=viewed_product_ids,
                skip=skip,
                take=take
            )

            # Fetch product details as gRPC objects
            grpc_products = _fetch_grpc_products(recommended_product_ids)
            return recommender_pb2.RecommendationResponse(
                recommended_products=grpc_products
            )

        except Exception as e:
            _handle_exception(context, f"Failed to get recommendations: {str(e)}")
            return recommender_pb2.RecommendationResponse()
        

def serve():
    # Train the model on startup
    print("Training recommender model...")
    try:
        recommender.train()
        print("Recommender model trained successfully!")
    except Exception as e:
        print(f"Error training model: {e}")
    
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    recommender_pb2_grpc.add_RecommenderServiceServicer_to_server(RecommenderServiceServicer(), server)
    server.add_insecure_port("[::]:50051")
    server.start()
    server.wait_for_termination()

if __name__ == "__main__":
    serve()