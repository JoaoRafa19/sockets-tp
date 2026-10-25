package tcp

// newRoutes liga cada método do protocolo ao handler dele.
func (h *Handler) newRoutes() map[string]HandlerFunc {
	return map[string]HandlerFunc{
		MethodCreateEvent:       h.createEvent,
		MethodReserveTicket:     h.reserveTicket,
		MethodGetReservation:    h.getReservation,
		MethodCancelReservation: h.cancelReservation,
	}
}
