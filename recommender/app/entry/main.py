import grpc
from concurrent import futures
from generated.pb import recommender_pb2_grpc, recommender_pb2
from app.services.recommender import recommender
from app.db.models import Product
from app.db.session import get_db_session



class RecommenderServiceServicer(recommender_pb2_grpc.RecommenderServiceServicer):
    def GetRecommendations(self, request, context):
        user_id = request.user_id

        try:
            recommended_product_ids = recommender.recommend(user_id, top_n=5)

            with get_db_session() as session:
                products = (
                    session.query(Product)
                    .filter(Product.id.in_(recommended_product_ids))
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

                return recommender_pb2.RecommendationResponse(
                    recommended_products=grpc_products
                )
        except Exception as e:
            context.set_details(f"Error getting recommendations: {str(e)}")
            context.set_code(grpc.StatusCode.INTERNAL)
            return recommender_pb2.RecommendationResponse()
        

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    recommender_pb2_grpc.add_RecommenderServiceServicer_to_server(RecommenderServiceServicer(), server)
    server.add_insecure_port("[::]:50051")
    server.start()
    server.wait_for_termination()

if __name__ == "__main__":
    serve()