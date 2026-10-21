import api from "./api";

export async function getCustomers(search = "") {
  const response = await api.get("/customers/", {
    params: search ? { search } : undefined,
  });
  return response.data;
}

export async function createCustomer(customer) {
  const response = await api.post("/customers/", customer);
  return response.data;
}

export async function updateCustomer(id, customer) {
  const response = await api.put(`/customers/${id}`, customer);
  return response.data;
}

export function deleteCustomer(id) {
  return api.delete(`/customers/${id}`);
}
