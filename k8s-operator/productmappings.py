
from tkinter import W
import kopf
import logging


@kopf.on.create('dashboard.freshworkscorp.com', 'v1', 'productmappings')
def create_fn(spec, name, **kwargs):
    logging.info(f'Creating product mapping {name}')
    logging.info(f'Spec: {spec}')
    
    
@kopf.on.update('dashboard.freshworkscorp.com', 'v1', 'productmappings')
def update_fn(spec, name, **kwargs):
    logging.info(f'Updating product mapping {name}')
    logging.info(f'Spec: {spec}')
