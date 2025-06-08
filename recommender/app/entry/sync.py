from kafka import KafkaConsumer
import json
from app.db.models import Product, Interaction
from app.db.session import get_db_session
from config.settings import KAFKA_BOOTSTRAP_SERVERS
import time
import threading
import requests
import logging

# Set up logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


def sync_products():
    logger.info("Starting product sync consumer...")
    consumer = KafkaConsumer(
        "product_events", 
        bootstrap_servers=KAFKA_BOOTSTRAP_SERVERS,
        group_id="recommender_product_sync",
        auto_offset_reset="earliest",
        value_deserializer=lambda m: json.loads(m.decode('utf-8'))
    )
    logger.info("Product sync consumer connected to Kafka")
    
    for message in consumer:
        event = message.value
        logger.info(f"Received product event: {event}")
        with get_db_session() as session:
            if event["type"] in ["product_created", "product_updated"]:
                product_data = event["data"]
                logger.info(f"Processing product event: {event['type']} for product ID: {product_data['product_id']}")
                product = session.query(Product).filter_by(id=product_data["product_id"]).first()
                if product:
                    product.name = product_data["name"]
                    product.description = product_data["description"]
                    product.price = product_data["price"]
                    product.account_id = product_data["accountID"]
                else:
                    product = Product(
                        id=product_data["product_id"],
                        name=product_data["name"],
                        description=product_data["description"],
                        price=product_data["price"],
                        account_id=product_data["accountID"]
                    )
                    session.add(product)
                session.commit()
                logger.info(f"Successfully synced product {product_data['product_id']}")
            elif event["type"] == "product_deleted":
                product = session.query(Product).filter_by(id=event["data"]["product_id"]).first()
                if product:
                    session.delete(product)
                    session.commit()
                    logger.info(f"Successfully deleted product {event['data']['product_id']}")


def fetch_product_from_graphql(product_id):
    """Fetch product from GraphQL API"""
    try:
        graphql_url = "http://graphql:8080/query"
        query = {
            "query": f'''
            query {{
                product(id: "{product_id}") {{
                    id
                    name
                    description
                    price
                    accountID
                }}
            }}
            '''
        }
        
        response = requests.post(graphql_url, json=query, timeout=5)
        response.raise_for_status()
        
        data = response.json()
        if "errors" in data:
            logger.error(f"GraphQL error fetching product {product_id}: {data['errors']}")
            return None
            
        products = data.get("data", {}).get("product", [])
        if products and len(products) > 0:
            product_data = products[0]
            return {
                "id": product_data["id"],
                "name": product_data["name"],
                "description": product_data["description"],
                "price": product_data["price"],
                "account_id": product_data["accountID"]
            }
        return None
    except requests.RequestException as e:
        logger.error(f"Failed to fetch product {product_id} from GraphQL: {e}")
        return None


def process_interactions():
    logger.info("Starting interaction sync consumer...")
    consumer = KafkaConsumer(
        "interaction_events",
        bootstrap_servers=KAFKA_BOOTSTRAP_SERVERS,
        group_id="recommender_interaction_sync",
        auto_offset_reset="earliest",
        value_deserializer=lambda m: json.loads(m.decode('utf-8'))
    )
    logger.info("Interaction sync consumer connected to Kafka")
    
    for message in consumer:
        event = message.value
        logger.info(f"Received interaction event: {event}")
        
        # Only process actual user interaction events
        if event["type"] not in ["purchase", "view", "add_to_cart"]:
            logger.info(f"Skipping non-interaction event: {event['type']}")
            continue
            
        # Check if the event has required fields for interactions
        if "user_id" not in event.get("data", {}):
            logger.warning(f"Skipping event without user_id: {event}")
            continue
            
        with get_db_session() as session:
            # Check if product exists in our database
            product = session.query(Product).filter_by(id=event["data"]["product_id"]).first()
            if not product:
                logger.info(f"Product {event['data']['product_id']} not found in database, trying to fetch from GraphQL...")
                # Try to fetch from GraphQL API
                product_data = fetch_product_from_graphql(event["data"]["product_id"])
                if product_data:
                    product = Product(
                        id=product_data["id"],
                        name=product_data["name"],
                        description=product_data["description"],
                        price=product_data["price"],
                        account_id=product_data["account_id"]
                    )
                    session.add(product)
                    session.commit()
                    logger.info(f"Successfully fetched and saved product {event['data']['product_id']} from GraphQL")
                else:
                    logger.warning(f"Failed to fetch product {event['data']['product_id']}, skipping interaction")
                    continue
            
            interaction = Interaction(
                user_id=event["data"]["user_id"],
                product_id=event["data"]["product_id"],
                interaction_type=event["type"]
            )
            session.add(interaction)
            session.commit()
            logger.info(f"Successfully processed interaction for product {event['data']['product_id']}")


if __name__ == "__main__":
    logger.info("Starting sync processes in parallel...")
    
    # Start both consumers in parallel threads
    product_thread = threading.Thread(target=sync_products, daemon=True)
    interaction_thread = threading.Thread(target=process_interactions, daemon=True)
    
    product_thread.start()
    interaction_thread.start()
    
    logger.info("Both sync processes started. Waiting...")
    
    # Keep main thread alive
    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        logger.info("Shutting down...")