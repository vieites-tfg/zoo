describe('Main page - Animals list', () => {
  beforeEach(() => {
    cy.intercept('GET', 'http://localhost:3000/animals', { fixture: 'animals.json' }).as('getAnimals')
    cy.visit('/')
    cy.wait('@getAnimals')
  })

  it('show table with animals', () => {
    cy.get('table').should('exist')
    cy.get('table tbody tr').should('have.length.greaterThan', 0)
    cy.get('table').contains('Leo').should('exist')
  })
})

